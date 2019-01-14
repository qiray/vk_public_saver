package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func parseJSON() {
	// TODO: different types
	type mmn struct {
		Rows []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"rows"`
	}

	body1 := []byte(`{
		"rows": [
			{
				"id": "01ae6145-90a3-11e7-7a69-8f55000cda4b",
				"type": "url",
				"url": "google.ru",
				"modificationsCount": 0
			},
			{
				"id": "11",
				"type": "image",
				"src": "image.png",
				"quantity": 1
			}
		]
	}`)
	var app = mmn{}
	err1 := json.Unmarshal(body1, &app)
	if err1 != nil {
		log.Fatal("error")
	}

	// итерируем
	for _, row := range app.Rows {
		fmt.Println(row.ID, row.Type)
	}
}
