// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("business_report"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Check if test user already exists
	var accounts []map[string]interface{}
	db.Table("wechat_accounts").Where("open_id = ?", "test").Find(&accounts)
	fmt.Printf("Existing test accounts: %d\n", len(accounts))
	for _, a := range accounts {
		fmt.Printf("  ID=%v OpenID=%v BusinessID=%v Phone=%v\n",
			a["id"], a["open_id"], a["business_id"], a["phone"])
	}

	// If not bound, create a binding for test user to business 2
	if len(accounts) == 0 {
		result := db.Exec(
			"INSERT INTO wechat_accounts (open_id, business_id, real_name, phone, status, created_at) VALUES (?, ?, ?, ?, ?, datetime('now'))",
			"test", 2, "测试用户", "13501126677", "active")
		if result.Error != nil {
			log.Fatal("Insert failed:", result.Error)
		}
		fmt.Println("Test user bound to business 2!")
	}
}
