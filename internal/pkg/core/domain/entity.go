package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Entity struct {
	id         uuid.UUID
	entityType string
	createdAt  time.Time
	updatedAt  time.Time
}

type EntityDataModel struct {
	Id         uuid.UUID `json:"id"          bson:"id,omitempty"`
	EntityType string    `json:"entity_type" bson:"entity_type,omitempty"`
	CreatedAt  time.Time `json:"created_at"  bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"  bson:"updated_at"`
}

type IEntity interface {
	Id() uuid.UUID
	CreatedAt() time.Time
	UpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
	SetEntityType(entityType string)
	SetId(id uuid.UUID)
}

// NewEntityWithId creates a new Entity with an id
func NewEntityWithId(id uuid.UUID, entityType string) *Entity {
	return &Entity{
		id:         id,
		createdAt:  time.Now(),
		entityType: entityType,
	}
}

// NewEntity creates a new Entity
func NewEntity(entityType string) *Entity {
	return &Entity{
		createdAt:  time.Now(),
		entityType: entityType,
	}
}

func (e *Entity) Id() uuid.UUID {
	return e.id
}

func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Entity) EntityType() string {
	return e.entityType
}

func (e *Entity) SetUpdatedAt(updatedAt time.Time) {
	e.updatedAt = updatedAt
}

func (e *Entity) SetEntityType(entityType string) {
	e.entityType = entityType
}

func (e *Entity) SetId(id uuid.UUID) {
	e.id = id
}
