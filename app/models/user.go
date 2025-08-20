package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id        uint   `gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique;type:varchar(20)" validate:"required,min=6,max=20"`
	Password  string `json:"password,omitempty" gorm:"type:varchar(255);" validate:"required,min=6"`
	FullName  string `json:"full_name" gorm:"type:varchar(100);" validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (i User) Validate() error {
	v := validator.New()
	return v.Struct(i)
}

type UserSession struct {
	Id                  uint `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UserId              uint      `json:"user_id" gorm:"type:int" validate:"required"`
	Token               string    `json:"token" gorm:"type:varchar(255)" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"type:varchar(255)" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
}

func (i UserSession) Validate() error {
	v := validator.New()
	return v.Struct(i)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (i LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(i)
}

type LoginResponse struct {
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (i LoginResponse) Validate() error {
	v := validator.New()
	return v.Struct(i)
}
