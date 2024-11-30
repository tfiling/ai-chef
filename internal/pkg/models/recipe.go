package models

type Recipe struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" validate:"required"`
}
