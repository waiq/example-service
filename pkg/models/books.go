package models

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Book struct {
	gorm.Model

	UUID   uuid.UUID
	Title  string
	Author string
}
