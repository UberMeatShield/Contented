package models

import (
	"time"
	"encoding/json"
    //"contented/actions"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
    "github.com/gobuffalo/nulls"
)

// MediaContainer is used by pop to map your media_containers database table to your go code.
type MediaContainer struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Src       string    `json:"src" db:"src"`
	Type      string    `json:"type" db:"type"`
	Preview   string    `json:"preview" db:"preview"`
    ContainerID nulls.UUID `json:"container_id" db:"container_id"`

    // Add this to the Container as well.
    // idx
    // is_active
}


// String is not required by pop and may be deleted
func (m MediaContainer) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MediaContainers is not required by pop and may be deleted
type MediaContainers []MediaContainer
type MediaMap map[uuid.UUID]MediaContainer

// String is not required by pop and may be deleted
func (m MediaContainers) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MediaContainer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MediaContainer) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MediaContainer) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
