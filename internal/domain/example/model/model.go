// Package model holds the domain entities for the example domain.
package model

import "time"

// Example is a minimal domain entity, kept deliberately trivial: its only job
// is to be a working end-to-end pattern a new service can copy (handler ->
// service -> repository -> MySQL) and then delete.
type Example struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateExampleRequest is the payload to create an Example.
type CreateExampleRequest struct {
	Name string `json:"name" validate:"required"`
}

// SetID sets the unique identifier for Example.
func (e *Example) SetID(id string) {
	e.ID = id
}

// SetCreatedAt sets the creation timestamp for Example.
func (e *Example) SetCreatedAt(createdAt time.Time) {
	e.CreatedAt = createdAt
}
