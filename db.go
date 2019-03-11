package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func dbExample() {
	db, err := sql.Open("sqlite3", "./test.db")
	defer db.Close()
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

func closeDatabase(db *sql.DB) {
	db.Close()
}

func createTable(db *sql.DB, initstring string) {
	stmt, err := db.Prepare(initstring)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

func initDataBase(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	checkErr(err)

	createTable(db,
		`CREATE TABLE IF NOT EXISTS posts (
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
		);`)

	createTable(db,
		`CREATE TABLE IF NOT EXISTS attachments (
			type TEXT,
			id INTEGER,
			owner_id INTEGER,
			post_id INTEGER,
			url TEXT,
			additional_info text,
			PRIMARY KEY (id, type, post_id)
		);`)

	return db
}

func savePosts(db *sql.DB, items []Post) {
	if len(items) == 0 {
		return
	}
	insertposts := `
		INSERT OR IGNORE INTO posts (
			id,
			from_id,
			owner_id,
			signer_id,
			date,
			marked_as_ads,
			post_type,
			text,
			is_pinned,
			comments_count,
			likes_count,
			reposts_count,
			views_count
		) VALUES 
	`
	insertattachmentsTemplate := `
		INSERT OR IGNORE INTO attachments (
			type,
			id,
			post_id,
			url,
			additional_info
		) VALUES 
	`

	insertattachments := insertattachmentsTemplate
	postsvalues := []interface{}{}
	attachmentsvalues := []interface{}{}
	count := 0

	tx, err := db.Begin() //start transaction
	checkErr(err)
	for _, item := range items {
		insertposts += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		postsvalues = append(postsvalues, item.ID, item.FromID, item.OwnerID, item.SignerID,
			item.Date, item.MarkedAsAds, item.PostType, item.Text, item.IsPinned,
			item.Comments.Count, item.Likes.Count, item.Reposts.Count, item.Views.Count)
		if len(item.Attachments) > 0 {
			for _, attachment := range item.Attachments {
				count++
				insertattachments += "(?, ?, ?, ?, ?),"
				if attachment.Type == "photo" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						attachment.Photo.ID, item.ID, "photo"+
							string(attachment.Photo.OwnerID)+"_"+string(attachment.Photo.ID),
						attachment.Photo.Text)
				} else if attachment.Type == "posted_photo" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						attachment.PostedPhoto.ID, item.ID, attachment.PostedPhoto.Photo604, "")
				} else if attachment.Type == "video" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						attachment.Video.ID, item.ID, "photo"+
							string(attachment.Video.OwnerID)+"_"+string(attachment.Video.ID),
						attachment.Video.Title)
				} else if attachment.Type == "audio" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						attachment.Audio.ID, item.ID, attachment.Audio.URL,
						attachment.Audio.Artist+"-"+attachment.Audio.Title)
				} else if attachment.Type == "doc" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						attachment.Doc.ID, item.ID, attachment.Doc.URL, attachment.Doc.Title)
				} else {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						item.ID, item.ID, "", "")
					//TODO: add other type
				}

				if count >= 500 {
					execInserts(tx, insertattachments, attachmentsvalues)
					count = 0
					insertattachments = insertattachmentsTemplate
					attachmentsvalues = []interface{}{}
				}
			}

		}
	}

	execInserts(tx, insertposts, postsvalues)
	if count > 0 {
		execInserts(tx, insertattachments, attachmentsvalues)
	}
	checkErr(tx.Commit()) //commit transaction
}

func execInserts(tx *sql.Tx, insertString string, values []interface{}) {
	insertString = strings.TrimSuffix(insertString, ",") //trim the last comma
	stmt, err := tx.Prepare(insertString)                //prepare the statement
	checkErr(err)
	_, err = stmt.Exec(values...) //format all values at once
	checkErr(err)
}

func savePostsResponse(db *sql.DB, p PostsResponse) {
	for _, val := range p.Response {
		savePosts(db, val.Items)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
