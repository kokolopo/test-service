package domain

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id int64) (*User, error)
}

type UserUsecase interface {
	CreateUser(name, email string) error
	GetAllUsers() ([]User, error)
	GetUserByID(id int64) (*User, error)
}
