// +build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
)

func main() {
	db, err := gorm.Open(sqlite.Open("business_report"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var businesses []map[string]interface{}
	db.Table("businesses").Find(&businesses)
	fmt.Println("=== Businesses ===")
	for _, b := range businesses {
		fmt.Printf("ID=%v Name=%v Phone=%v\n", b["id"], b["name"], b["phone"])
	}

	var users []map[string]interface{}
	db.Table("wechat_accounts").Find(&users)
	fmt.Println("=== WeChat Accounts ===")
	for _, u := range users {
		fmt.Printf("ID=%v OpenID=%v BusinessID=%v Phone=%v\n", u["id"], u["openid"], u["business_id"], u["phone"])
	}
}
