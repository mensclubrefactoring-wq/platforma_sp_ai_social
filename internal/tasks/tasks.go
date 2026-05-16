package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"platforma-sp/internal/db"
	"platforma-sp/internal/shared"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if db.DB == nil {
		json.NewEncoder(w).Encode([]shared.Task{})
		return
	}
	var tasks []shared.Task
	query := db.DB.Model(&shared.Task{})

	userIDStr := r.Header.Get("User-ID")
	var userID uint
	if userIDStr != "" {
		fmt.Sscanf(userIDStr, "%d", &userID)
	}

	fmt.Printf("DEBUG: GetTasksHandler. UserID: %d, Role: %s\n", userID, r.Header.Get("User-Role"))

	if userID != 0 {
		// Каждая роль видит только свои задачи согласно требованию
		query = query.Where("creator_id = ?", userID)
	} else {
		// Если не залогинен, возвращаем пустой список (безопасность)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]shared.Task{})
		return
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	category := r.URL.Query().Get("category")
	if category != "" && category != "Все" {
		query = query.Where("category = ?", category)
	}
	location := r.URL.Query().Get("location")
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	query.Order("created_at desc").Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if db.DB == nil {
		http.Error(w, "Database not configured", 500)
		return
	}
	var task shared.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	userIDStr := r.Header.Get("User-ID")
	var userID uint
	fmt.Sscanf(userIDStr, "%d", &userID)
	task.CreatorID = userID
	task.Status = "active"

	if err := db.DB.Create(&task).Error; err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(task)
}
