package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func dbExample() {
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	initstring := `CREATE TABLE IF NOT EXISTS userinfo (
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NULL,
        departname VARCHAR(64) NULL,
        created DATE NULL
	);`

	stmt, err := db.Prepare(initstring)
	checkErr(err)
	stmt.Exec()

	// insert
	stmt, err = db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	// update
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	rows.Close() //good habit to close

	// delete
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func initDataBase(filepath string) {
	db, err := sql.Open("sqlite3", filepath)
	checkErr(err)

	initstring := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER,
			from_id INTEGER,
			owner_id INTEGER,
			signer_id INTEGER,
			date INTEGER,
			marked_as_ads INTEGER,
			post_type TEXT,
			text TEXT,
			is_pinned INTEGER,
			comments_count INTEGER,
			likes_count INTEGER,
			reposts_count INTEGER,
			views_count INTEGER,
			PRIMARY KEY (id, from_id)
		);
	`
	//TODO: add attachments

	stmt, err := db.Prepare(initstring)
	checkErr(err)
	stmt.Exec()
}

func savePost(filepath string) { //https://stackoverflow.com/questions/21108084/golang-mysql-insert-multiple-data-at-once

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
