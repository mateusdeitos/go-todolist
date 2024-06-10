package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusdeitos/go-todolist/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := mux.NewRouter()
	db := createDb()
	db.AutoMigrate(&entity.TodoEntity{})

	r.HandleFunc("/todo", ListTodosHandler).Methods(http.MethodGet)

	http.ListenAndServe(":9000", r)
}

func createDb() *gorm.DB {
	dsn := "host=localhost user=admin password=123 dbname=gorm port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos := []*entity.TodoEntity{
		entity.NewTodo("todo1"),
		entity.NewTodo("todo2"),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(todos)
}
