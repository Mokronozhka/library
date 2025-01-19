package models

import (
	"github.com/google/uuid"
	"time"
)

type UserStruct struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name" validate:"required"`
	Password         string    `json:"pwd" validate:"required,min=8"`
	Email            string    `json:"email" validate:"required,email"`
	Age              int       `json:"age,omitempty" validate:"gte=14"`
	DateRegistration time.Time `json:"date_reg,omitempty"`
}

type UserLoginStruct struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"pwd" validate:"required,min=8"`
}

type BookStruct struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"desc,omitempty"`
	Author      string    `json:"author" validate:"required"`
	DateWriting time.Time `json:"date_wrt,omitempty"`
}
