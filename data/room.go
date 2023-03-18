package data

import (
	"context"
	"strings"
	"time"
)

// Room represent rooms table in the database
type Room struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Room) Table() string {
	return "rooms"
}

// Create inserts a room into the database
func (r *Room) Create(room Room) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	room.Name = strings.ToLower(room.Name)

	var newID int
	query := `insert into rooms (name, created_at, updated_at) values($1, $2, $3) returning id`
	err := DB.QueryRowContext(ctx, query,
		room.Name,
		time.Now(),
		time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// GetAll returns all rooms from the database
func (r *Room) GetAll() ([]Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []Room

	rows, err := DB.QueryContext(ctx, "SELECT * FROM rooms ORDER BY name")
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetById returns a room by id
func (r *Room) GetById(id int) (Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room Room

	query := ` select id, name, created_at, updated_at from rooms where id=$1`

	row := DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.Name,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

// GetByName returns a room by name
func (r *Room) GetByName(name string) (Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name = strings.ToLower(name)

	var room Room

	query := ` select id, room_name, created_at, updated_at from rooms where name=$1`

	row := DB.QueryRowContext(ctx, query, name)
	err := row.Scan(
		&room.ID,
		&room.Name,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

// IsAvailable checks if a room is available for a given time period
//
// if the desired range does not overlap with any restriction, the room is available
func (r *Room) IsAvailable(roomID int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	query := `
			select 
				count(id) 
			from 
				restrictions
			where 
			      room_id = $1
				and ($2 < end_date and $3 > start_date)`

	row := DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	//no rows found, means room is free
	if count == 0 {
		return true, nil
	}
	return false, nil
}

// GetAnyAvailable returns zero or more rooms that are available for a given time period
func (r *Room) GetAnyAvailable(start, end time.Time) ([]Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []Room

	query := `
		select 
			r.id, r.name
		from
			rooms r
		where r.id not in 
			(select 
				rest.room_id 
			from 
				restrictions rest 
			where 
			$1 < rest.end_date and $2 > rest.start_date)`

	rows, err := DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room

		err := rows.Scan(
			&room.ID,
			&room.Name,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}
