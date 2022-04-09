package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init(){
	//DB initialization
	db, err := sql.Open("mysql", "root:appliCATION123@#@tcp(127.0.0.1:3306)/new_DB")
	if err != nil {
		fmt.Println("error validating sql.Open arguments")
		panic(err.Error())
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("error verifying connection with db.Ping")
		panic(err.Error())
	}

	DB = db
	fmt.Println("successfully connected to database!")

	//--------initial inserting of my populated data into the DB----------------------
	//  query, err  := DB.Prepare("INSERT INTO BLOG (id, title, ingredients, content, time, author_name) VALUES (?, ?, ?, ?, ?, ?)")
	//    if err != nil {
	// 	   log.Printf("%s", err)
	// 	   return
	//   }
	//   for _, b := range blog {
	// 	  _, err = query.Exec(b.ID, b.Title, b.Ingredients, b.Content, b.Time, b.AuthorName)
	// 	  if err != nil {
	// 		  log.Printf("%s", err)
	// 		  return
	// 	  }
	//   }

}