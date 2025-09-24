package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct{
	DB *sql.DB
}

// Just update this struct - use a string for the date field
type Event struct{
	Id          int       `json:"id"`
	OwnerId     int       `json:"ownerId"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        string    `json:"date" binding:"required"`
	Location    string    `json:"location" binding:"required,min=3"`
}

// Update SQL placeholder to ? for SQLite
func (m *EventModel) Insert(event *Event) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO events (owner_id, name, description, date, location)
	VALUES(?, ?, ?, ?, ?)`
	
	// Use Exec + LastInsertId for SQLite
	res, err := m.DB.ExecContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location)
	if err != nil {
		return err
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	
	event.Id = int(id)
	return nil
}

// Also update remaining methods to use ? placeholders
func (m *EventModel) GetAll() ([]*Event, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil{
		return nil,err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next(){
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil{
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return events, nil
}

// GetPage retrieves a specific page of events with pagination
func (m *EventModel) GetPage(page int, pageSize int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	query := `SELECT id, owner_id, name, description, date, location FROM events LIMIT ? OFFSET ?`

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil{
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next(){
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil{
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return events, nil
}

func (m *EventModel) Get(id int) (*Event, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events WHERE id = ?`

	var event Event
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&event.Id,
		&event.OwnerId,
		&event.Name,
		&event.Description,
		&event.Date,
		&event.Location,
	)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func (m *EventModel) Update(event *Event) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE events SET name=?, description=?, date=?, location=? WHERE id=?`
	_, err := m.DB.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.Id)
	if err != nil{
		return err
	}

	return nil
}

func (m *EventModel) Delete(id int) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `DELETE FROM events WHERE id=?`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil{
		return err
	}

	return nil
}

// GetCount returns the total number of events
func (m *EventModel) GetCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `SELECT COUNT(*) FROM events`
	
	var count int
	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}