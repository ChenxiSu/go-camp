package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	errors "github.com/pkg/errors"

	"hm1/service"
)

func main() {
	err := seed_db()
	if err != nil {
		panic(err)
	}

	user, err := service.GetUserById(2)
	if err != nil {
		fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("\nstack trace: %+v\n", err)
	}else {
		fmt.Println(user)
	}
}

func seed_db() error {
	file, err := os.Create("sqlite_database.db") // Create SQLite file
	if err != nil {
		fmt.Println(err.Error())
	}
	file.Close()

	db, err := sql.Open("sqlite3", "./sqlite_database.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close() 

	// create table
	create_table_sql := `
		CREATE TABLE userinfo (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"username" VARCHAR(64) NULL,
			"email" DATE NULL
		);
	`
	statement, err := db.Prepare(create_table_sql) // Prepare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, email) values(?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec("morningsu", "test@gmail.com")
	if err != nil {
		return err
	}
	return nil
}