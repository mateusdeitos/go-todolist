package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type TodoEntity struct {
	*gorm.Model
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

func NewTodo(name string) *TodoEntity {
	return &TodoEntity{
		Name:     name,
		Complete: false,
	}
}

func (t *TodoEntity) ToggleComplete() {
	t.Complete = !t.Complete
}

func (t *TodoEntity) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":       t.ID,
		"name":     t.Name,
		"complete": t.Complete,
	})
}
