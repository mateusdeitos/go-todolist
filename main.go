package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/mateusdeitos/go-todolist/src/entity"
)

func main() {
	fmt.Println("Server started")
	ctx := context.Background()
	r := mux.NewRouter()
	conn, err := getDBConn(ctx)
	if err != nil {
		fmt.Println("Error connecting to db")
		log.Fatal(err)
		return
	}

	fmt.Println("Connected to database")
	defer conn.Close(ctx)

	queries := entity.New(conn)

	r.HandleFunc("/todo", Wrapped(ctx, queries, ListTodosHandler)).Methods(http.MethodGet)
	r.HandleFunc("/todo", Wrapped(ctx, queries, CreateTodoHandler)).Methods(http.MethodPost)

	r.HandleFunc("/todo/{id:[0-9]+}", Wrapped(ctx, queries, UpdateTodoHandler)).Methods(http.MethodPut)

	done := make(chan os.Signal, 1)

	go func() {
		err := http.ListenAndServe(":9000", r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Server running on port 9000")
	}()

	<-done

	fmt.Println("Server stopped")
}

func getDBConn(ctx context.Context) (*pgx.Conn, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	return pgx.Connect(ctx, dsn)
}

type CustomHandler func(context.Context, http.ResponseWriter, *http.Request, *entity.Queries)

func Wrapped(ctx context.Context, q *entity.Queries, h CustomHandler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h(ctx, w, r, q)
	}

	return http.HandlerFunc(fn)
}

type PaginatedResult[T any] struct {
	Items  []T   `json:"items"`
	Count  int64 `json:"count"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func ListTodosHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, q *entity.Queries) {
	var limit int
	var offset int
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	result, err := q.ListTodos(ctx, entity.ListTodosParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	count, err := q.CountTodos(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(PaginatedResult[entity.Todo]{
		Items:  result,
		Count:  count,
		Limit:  limit,
		Offset: offset,
	})
}

func CreateTodoHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, q *entity.Queries) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var todo entity.Todo
	err := decoder.Decode(&todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err = q.CreateTodo(ctx, entity.CreateTodoParams{
		Name:     todo.Name,
		Complete: todo.Complete,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodoHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, q *entity.Queries) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := q.GetTodo(ctx, int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = q.UpdateTodo(ctx, entity.UpdateTodoParams{
		Name:     todo.Name,
		Complete: todo.Complete,
		ID:       int64(id),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func GetTodoHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, q *entity.Queries) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := q.GetTodo(ctx, int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodoHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, q *entity.Queries) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = q.DeleteTodo(ctx, int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
