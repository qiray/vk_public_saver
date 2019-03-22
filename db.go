// vk_public_saver - tool for saving walls data from vk.com
// Copyright (c) 2019, Yaroslav Zotov, https://github.com/qiray/
// All rights reserved.

// This file is part of vk_public_saver.

// vk_public_saver is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// vk_public_saver is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with vk_public_saver.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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
			for i, attachment := range item.Attachments {
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
				} else if attachment.Type == "link" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						i, item.ID, attachment.Link.URL, attachment.Link.Title+":"+
							attachment.Link.Caption)
				} else if attachment.Type == "poll" {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						attachment.Poll.ID, item.ID, "wall"+string(attachment.Poll.OwnerID)+"_"+
							string(item.ID), attachment.Poll.Question)
				} else {
					attachmentsvalues = append(attachmentsvalues, attachment.Type,
						item.ID, item.ID, "", "")
					// fmt.Println(attachment.Type)
					//Add other type (maybe)
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
