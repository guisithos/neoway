package models

import (
	"time"
)

type ClientType string

const (
	PersonType   ClientType = "PERSON"
	BusinessType ClientType = "BUSINESS"
)

type Client struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null"`
	Document  string     `json:"document" gorm:"unique;not null"` // CPF ou CNPJ
	Type      ClientType `json:"type" gorm:"not null"`
	Blocked   bool       `json:"blocked" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
