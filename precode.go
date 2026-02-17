package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func GetTasks(w http.ResponseWriter, r *http.Request) {

	ans, err := json.Marshal(tasks)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(ans)
}

func PostTasks(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    tasks[newTask.ID] = newTask

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	
	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(w, "no found", http.StatusBadRequest)
	}
	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task := tasks[id]

	convTask, err := json.Marshal(task)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(convTask)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", GetTasks)
	r.Post("/tasks", PostTasks)
	r.Delete("/tasks/{id}", DeleteTask)
	r.Get("/tasks/{id}", GetTask)

	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
