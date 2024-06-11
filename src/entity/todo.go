package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Todo struct {
	*gorm.Model
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

func NewTodo(name string) *Todo {
	return &Todo{
		Model:    &gorm.Model{},
		Name:     name,
		Complete: false,
	}
}

func (t *Todo) ToggleComplete() {
	t.Complete = !t.Complete
}

func (t *Todo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":       t.ID,
		"name":     t.Name,
		"complete": t.Complete,
	})
}
