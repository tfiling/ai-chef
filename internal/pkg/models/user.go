package models

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"name" validate:"required"`
}
