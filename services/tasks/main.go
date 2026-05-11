package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"platforma-sp/internal/shared"
)

type DB struct {
	Tasks []shared.Task `json:"tasks"`
}

func getDB() DB {
	var db DB
	data, _ := os.ReadFile("db_tasks.json")
	json.Unmarshal(data, &db)
	return db
}

func saveDB(db DB) {
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile("db_tasks.json", data, 0644)
}

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(getDB().Tasks)
			return
		}
		if r.Method == "POST" {
			var task shared.Task
			json.NewDecoder(r.Body).Decode(&task)
			task.ID = fmt.Sprintf("%d", time.Now().UnixNano())
			task.CreatedAt = time.Now()
			task.Status = "active"
			
			db := getDB()
			db.Tasks = append(db.Tasks, task)
			saveDB(db)
			json.NewEncoder(w).Encode(task)
			return
		}
	})

	fmt.Println("Tasks Service running on :3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
