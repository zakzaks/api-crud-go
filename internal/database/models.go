package database

import "database/sql"

type Models struct {
	// Add your model structs here, e.g. User, Product, etc.
	User     UserModel
	Events  EventModel
	Attendees AttendeeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User:     UserModel{DB: db},
		Events:  EventModel{DB: db},
		Attendees: AttendeeModel{DB: db},
	}
}