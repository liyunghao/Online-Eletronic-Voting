# Leader and Followers Cluster

Also known as: `Primary Backup` , `Passive Replication` , `Master-Slave Pattern`

設計這個 System 時有以下幾個要考量的面相：

1. Cluster Management (Creation, Control)
2. Leader & Followers Communication Pattern (Heartbeat, Write-Sync)
3. Leader Failed (Delcared Dead, New Leader)
4. Node Failed (Catch Up Process, Re-Join the Cluster)
5. Hardening Storage System (Replication Logs)
6. Problem Avoidance (Split Brain, Conflict after Agreement on New Leader)

基本設定：

- 管理 API 將以 RESTful API 來實作，簡化純 TCP 實作上的問題
  - 所有路徑均以 `POST` 方法實作
- 管理 Port: 9000
  - (暫時不管資安上的問題，未來可以掛上 SSL 以及加上 Cluster 認證來處理這個問題)
- 目前的設計非常不好 Scale，未來有機會可以修正
  - Sync Communication
  - Pass All Data Conflict Solving

## Cluster Mangement

中心設計：**Leader 管理整個 Cluster**

每個 Node 都要去吃一個 cluster topology 的 json configuration，格式如下：

```json
{
  // Self-node configuration
  "node": {
    "name": "node-1",
    "id": 1
  },
  // Default configuration. When the node join the cluster, it need to acquire
  // newest state from leader
  "clusters": [
    {
      "name": "node-1",
      "ip": "192.168.1.10",
      "id": 1
    }
    // .....
  ]
}
```

### Cluster Creation / Node Join Existing Cluster

所有 Node 一開始進入的狀態會直接以設定檔的模式開始運作，但在真正上線開始前會先經過以下 Handshake
去確認是否自己現在的 Role 是合適的：

1. 向 Topology 中所有節點確認是否有 **Leader** 存在
2. 如果 Leader 不存在則依據 [Leader 產生機制](#new-leader)開始產生新的 Leader
3. 如果 Leader 存在則開始執行 Follower Role，呼叫 Catch Up Process，開始接受 HeartBeat
   Message

### Control Logic

Cluster Status Update - Node Index：

目前先寫死，數量固定的 Node

## Leader & Followers Communication Pattern

這裡重點是 Write-Sync 要怎麼實作。目前因為我們 Cluster Node 並不會太多，所以全數採用 Sync 的
溝通模式，bad scalablity but good consistency。

### Heartbeat

- Route: `/hearbeat`
- Return: nothing, `code: 200`

每隔 `30 秒` Leader 會去 Poll 每個 Follower，確認自己還活著。

### Write Sync

- Route: `/writesync`
- Correct Return: nothing, `code: 200`
- Fail Return: error message, `code: 對應 Status Code`
- Leader 的 `/writesync` 並不會運作

資料以 `json` 格式來傳輸，格式舉例如下：

```json
{
  "storage_cmd": "CreateUser",
  // Payload 格式均已 Storage 那邊定義的 Field Tag 為主，直接 Unmarshal
  "payload": "JSON 格式的 Payload"
}
```

也許有更 Compact 的傳輸格式，以利減少整體系統架構上的 Overhead 或效能耗損，但目前先以這種格式以簡化
Implementation 上的困難。

## Leader Failed

### Declared Dead

當 Follower 在收到上一個 Heartbeat 後的 `1 min` 內沒有收到上一個 Message 則開始進行
[Leader 產生機制](#new-leader)

### New Leader

- Route: `/declare_capability`
- Payload: `{ "node_idx": [node id] }`

採用的方法很簡單，因為每個 Node 有編號，就照著 Index 順序去確認下一個 Leader 是誰，數字越小
Priority 越高。演算法如下：

1. 每個活著的 Node 都可以 Broadcast 自己成為 Leader 的意願
   - 每 `5 秒` Broadcast 一次
2. Timeout `30 秒`
3. 根據 Priority 產生新的 Leader
   - 自己是最大 Priority 則直接成為 Leader
   - 其他人則成為 Followers

這樣設計難免會有 Split Brain 的問題，所以等等要處理 Split Brain 的問題

### Election of New Leader

- Route: `/recv_elect`
- Payload: `{ "approved": [bool]}`

接收別人傳來的 declare_capability, 回傳是否同意新的 Leader
(Payload 跟機制可能可以再改，只是先定義好 api route)

## Node Failed

### Catch up process

- Route: `/catch_up`
- Payload: `{ "snapshot_id": 1 }`
- Response: `{ logs: [ [很多 WriteSync 的 Log] ] }`
- Only leader will open this route

Follower 會把自己目前擁有最新的 Snapshot 的 Index 傳給 Leader，Leader 則會回傳最新的 Log
回來。這裡有可能會需要處理 **Conflict** 的問題。

## Hardening Storage System

In-Memory 的 Storage 可以透過紀錄 Logic Logs 的方式來增加 Reliablity，一方面 System
Crash 時有資料可以復原，另一方面是可以提供 Catch-up Process 資料來復原。

Log 格式長這樣：`[type_id]|[JSON string or Method Parameters as JSON]`

### Log Type 1: User Creation Log

- Format: `1|User JSON`

### Log Type 2: User Revoke Log

- Format: `2|RemoveUser Parameters as JSON`
- Example: `2|{"name": "tester"}`

### Log Type 3: Create Election

- Format: `3|Election JSON`

### Log Type 4: Vote Election

- Format: `4|VoteElection Parameters`
- Example: `4|{"election_name": "election1", "voter_name": "tester", "choice": "cand1"}`

### Snapshot 建立模式與時機

每 10 筆資料寫入後，需要創建一個 Snapshot。Follower 與 Leader 獨立自行建立 Snapshot，
而且是背景處理，但同時 Log 還要繼續新增

## Problem Avoidance

- Route: `/service_down`
- Route: `/service_up`

### Split Brain

- Route: `/is_leader`
- Response: `true or false`
- Leader Polling 其他 Node 的 API Route

Leader 每隔 `1 min` 會去 Polling 有沒有其他 Node 也是 Declare 自己是 Leader。當發現 Cluster
中有兩個或兩個以上的 Leader，透過以下方法來解決：

1. System Wide Shutdown
2. Priority 比較高（index 最小）的會成為唯一的 Leader
3. 其他將要被取代掉的 Leader 把目前 DB 中的資料送給新 Leader 處理 Conflict
4. Conflict 處理邏輯請參考[下一段文件](#conflict)
5. Leader Node 創建一個新的 Snapshot，刪除所有 Replication Logs
6. Followers 強制更新所有資料，產新的 Snapshot，刪除 Replication Logs
7. System Wide Up

### Conflict

Follower

- Route: `/override_all_data`
- Payload 就是那一大堆資料
- 當 Leader 處理完 Conflict 資料後拿來強迫 Follower 更新資料
- 創建一個 ID 為 0 的 Snapshot

Leader

- Route: `/merge_conflict`(Leader)
- Payload 是自己擁有的全部資料
- 將要被替換掉的 Leader 傳回自己的資料給最後的 Leader

處理 Conflict 的演算法：

1. 資料兩邊如果有聯集的直接新增
2. 如果兩邊資料有衝突，如 Vote 數量，等等資料上的不一樣，以 Priority 高的 Node 資料為主
