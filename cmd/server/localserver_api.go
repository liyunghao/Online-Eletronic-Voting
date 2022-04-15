package main

import (
	"fmt"

	db "github.com/liyunghao/Online-Eletronic-Voting/internal/server/database"
)

// Minimum Requirements
func RegisterVoter(name string, group string, public_key string) {
	_, err := db.SqliteDB.Exec("INSERT INTO voters (name, grouptype, public_key) VALUES (?, ?, ?)", name, group, public_key)
	if err != nil {
		fmt.Println(err)
	}
}

func UnregisterVoter(name string) {
	_, err := db.SqliteDB.Exec("DELETE FROM voters WHERE name = ?", name)
	if err != nil {
		fmt.Println(err)
	}
}

// --------------------------------------------------
// Miscellaneous APIs
