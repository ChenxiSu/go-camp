package dao

import (
	"database/sql"
	"fmt"

	errors "github.com/pkg/errors"

	"hm1/model"
)

func GetUserById(id int64) (*model.User, error) {
	db, err := sql.Open("sqlite3", "./sqlite_database.db")
	if err != nil {
		return nil, errors.New("Failed to connect to db")
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM userinfo where id = %d;", id)
	row := db.QueryRow(query)

	var userid int64
	var username string
	var email string
	
	err = row.Scan(&userid, &username, &email)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("No result found based on %d", id))
	}
	user := model.User{userid, username, email}
	return &user, nil
}