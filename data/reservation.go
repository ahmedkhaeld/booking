package data

import (
	"context"
	"github.com/ahmedkhaeld/jazz/forms"
	"strings"
	"time"
)

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
	Processed int
}

func (r *Reservation) Table() string {
	return "reservations"
}

func (r *Reservation) Validate(v *forms.Form) {
	v.Required("first_name", "last_name", "email", "phone")
	v.MinLength("first_name", 3)
	v.MinLength("last_name", 3)
	v.IsEmail("email")

}

func (r *Reservation) Create(res Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res.FirstName = strings.ToLower(res.FirstName)
	res.LastName = strings.ToLower(res.LastName)
	res.Email = strings.ToLower(res.Email)

	var newID int
	query := `insert into reservations (first_name, last_name, email, 
			phone, start_date, end_date, room_id, created_at, updated_at)
			values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := DB.QueryRowContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *Reservation) GetAll() ([]Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []Reservation

	query := `
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	order by r.start_date asc
`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.RoomID,
			&i.Room.Name,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (r *Reservation) GetByID(id int) (Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res Reservation
	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where r.id = $1 
	`

	row := DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.Name,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *Reservation) GetOnlyNew() ([]Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []Reservation

	query := `
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where processed = 0 
	order by r.start_date asc
`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.RoomID,
			&i.Room.Name,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (r *Reservation) Update(res Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5
		where id =$6
`

	_, err := DB.ExecContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		time.Now(),
		res.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reservation) UpdateProcessedStatus(processed, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "update reservations set processed = $1 where id =$2"

	_, err := DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil

}

func (r *Reservation) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "delete from reservations where id =$1"

	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
