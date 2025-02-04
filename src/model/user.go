package model

import "github.com/google/uuid"

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordRequest struct {
	Code     string `json:"code"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
}

type CreateUserRequest struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Phone     string    `json:"phone"`
	RoleID    uuid.UUID `json:"role_id"`
	Code      string    `json:"code"`
}
