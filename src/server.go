package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusdeitos/go-todolist/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := mux.NewRouter()
	db := createDb()
	err := db.AutoMigrate(&entity.Todo{})
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/todo", ListTodosHandler).Methods(http.MethodGet)

	http.ListenAndServe(":9000", r)

	fmt.Println("Server running on port 9000")
}

func createDb() *gorm.DB {
	dsn := "host=db user=admin password=123 dbname=go_todo_list port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos := []*entity.Todo{
		entity.NewTodo("todo1"),
		entity.NewTodo("todo2"),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(todos)
}
