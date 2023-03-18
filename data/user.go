package data

import (
	"context"
	"errors"
	"github.com/ahmedkhaeld/jazz"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) Table() string {
	return "users"
}

func (u *User) Validate(v *jazz.Validation) {
	v.Check(u.FirstName != "", "first_name", "First Name must be provided")
	v.Check(u.LastName != "", "last_name", "Last Name must be provided")
	v.Check(u.Password != "", "password", "Password must be provided")
	v.Check(u.Email != "", "email", "Email must be provided")
	v.IsEmail("email", u.Email)

}

func (u *User) Insert(user User) (int, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	hashStr := string(hash)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//forcefully store emails by lowercase
	user.Email = strings.ToLower(u.Email)

	query := `insert into users (first_name, last_name, email, 
			password, created_at, updated_at, access_level)
            values ($1, $2, $3, $4, $5, $6, $7)`

	result, err := DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		hashStr,
		time.Now(),
		time.Now(),
		user.AccessLevel,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}
func (u *User) InsertAndReturnId(user User) (int, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	hashStr := string(hash)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	query := `insert into users (first_name, last_name, email, 
			password, created_at, updated_at, access_level)
            values ($1, $2, $3, $4, $5, $6, $7)`

	//forcefully store emails by lowercase
	user.Email = strings.ToLower(u.Email)

	err = DB.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		hashStr,
		time.Now(),
		time.Now(),
		user.AccessLevel).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil

}

func (u *User) SelectAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var users []User

	rows, err := DB.QueryContext(ctx, "SELECT * FROM users ORDER BY last_name, first_name")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.AccessLevel,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (u *User) Count() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (u *User) GetByID(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at 
			from users where id=$1`

	row := DB.QueryRowContext(ctx, query, id)
	var user User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil

}
func (u *User) GetByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at 
			from users where email=$1`

	row := DB.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil

}

func (u *User) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`

	_, err := DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hash string

	email = strings.ToLower(email)

	row := DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email=$1", email)
	err := row.Scan(&id, &hash)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	}
	if err != nil {
		return 0, "", err
	}

	return id, hash, nil

}

func (u *User) UpdatePassword(user User, newHash []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hash := string(newHash)

	query := `update users set password=$1 where id=$2`
	_, err := DB.ExecContext(ctx, query, hash, user.ID)
	if err != nil {
		return err
	}
	return nil

}

func (u *User) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx, "DELETE FROM users WHRE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
