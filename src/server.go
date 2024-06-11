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

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDb.Close()

	r.HandleFunc("/todo", Wrapped(db, ListTodosHandler)).Methods(http.MethodGet)
	r.HandleFunc("/todo", Wrapped(db, CreateTodoHandler)).Methods(http.MethodPost)


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

type CustomHandler func(http.ResponseWriter, *http.Request, *gorm.DB)

func Wrapped(db *gorm.DB, h CustomHandler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h(w, r, db)
	}

	return http.HandlerFunc(fn)
}

func ListTodosHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	result := []entity.Todo{}
	todos := db.Find(&result)
	if todos.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var todo entity.Todo
	err := decoder.Decode(&todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Create(&todo)
	json.NewEncoder(w).Encode(todo)
}
}
