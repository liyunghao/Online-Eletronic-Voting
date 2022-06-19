# Active-Active 架構

基本上實作跟 [Leader Follower Cluster](./leader_follower_cluster.md) 中提到的實作方向一樣，
只是多了以下幾個部分，以符合 Active-Active 架構

## Same-tier nodes

Active-Active 架構中所有 Node 都是平等的，所以在原本 Leader Follower Cluster 架構下，把
原本限制只有 Primary 可以接受請求的限制拔掉，變成兩個 Node 都可以接受請求以及 Sync 資料到另一個 Node，
進而達到 High Available 的目的。

## Conflict Solving

因為是 Voting 這個行為的關係，當今天兩個 Node 看到的資料不相同時，兩邊 Merge 的實作可以很簡單的把
投過票的人合併在一起，重複的則不算，解決 Network Partition 下兩邊 Node 各自為政的問題
