package data

import (
	"context"
	"log"
	"time"
)

// Restriction represents a booking for a room for a given date range
type Restriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
}

func (r *Restriction) Table() string {
	return "restrictions"
}

// Create inserts a restriction into the database
// restrictions are used to block out dates when a room is not available
// for example, when a room is booked or being cleaned or maintained
func (r *Restriction) Create(restrict Restriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at)
             values ($1, $2, $3, $4, $5, $6)`

	_, err := DB.ExecContext(ctx, query,
		restrict.StartDate,
		restrict.EndDate,
		restrict.RoomID,
		restrict.ReservationID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil
	}
	return nil
}

// GetForRoom returns all restrictions for a given room
func (r *Restriction) GetForRoom(start, end time.Time, roomID int) ([]Restriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []Restriction

	query := ` 
		select id, reservation_id, room_id, start_date, end_date
		from restrictions 
		where $1 < end_date and $2 >= start_date and room_id = $3 
`
	rows, err := DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rest Restriction
		err := rows.Scan(
			&rest.ID,
			&rest.ReservationID,
			&rest.RoomID,
			&rest.StartDate,
			&rest.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, rest)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

func (r *Restriction) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from restrictions where id = $1`

	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
