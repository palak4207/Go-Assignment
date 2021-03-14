package db

import (
	"context"
	"time"

	"github.com/gtldhawalgandhi/go-training/3.Intermediate/util"
)

type UserRequest struct {
	UserName  string    `json:"user_name" binding:"required,alphanum"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6"`
	FullName  string    `json:"full_name"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// UserResponse ...
type UserResponse struct {
	UserName       string    `json:"user_name"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	HashedPassword string    `json:"pass_hash"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

// GetUserByUsername ..
func (pg *PGStore) GetUserByUserName(ctx context.Context, userName string) (UserResponse, error) {
	var ur UserResponse
	err := pg.db.QueryRow(context.Background(), "select user_name, first_name, last_name, created_at, pass_hash from users where user_name=$1", userName).Scan(&ur.UserName, &ur.FirstName, &ur.LastName, &ur.CreatedAt, &ur.HashedPassword)
	if err != nil {
		return UserResponse{}, err
	}

	return ur, nil
}

// GetUserByEmail ..
func (pg *PGStore) GetUserByEmail(ctx context.Context, email string) (UserResponse, error) {
	var ur UserResponse
	err := pg.db.QueryRow(context.Background(), "select user_name, first_name, last_name, created_at from users where email=$1", email).Scan(&ur.UserName, &ur.FirstName, &ur.LastName, &ur.CreatedAt)
	if err != nil {
		return UserResponse{}, err
	}

	return ur, nil
}

// CreateUser ..
func (pg *PGStore) CreateUser(ctx context.Context, user UserRequest) (UserResponse, error) {
	var ur UserResponse

	passHash, err := util.HashPassword(user.Password)

	if err != nil {
		return UserResponse{}, err
	}

	err = pg.db.QueryRow(context.Background(), `
	insert into users (user_name, first_name, last_name, email, pass_hash, created_at) values 
		($1,$2,$3,$4,$5,$6)
	on conflict (user_name) do 
		update set 
			user_name = excluded.user_name,
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			email = excluded.email,
			pass_hash = excluded.pass_hash
	RETURNING user_name;
	`, user.UserName, user.FirstName, user.LastName, user.Email, passHash, time.Now()).Scan(&ur.UserName)
	if err != nil {
		return UserResponse{}, err
	}

	return ur, nil
}

// UpdateUser ..
func (pg *PGStore) UpdateUser(ctx context.Context, user UserRequest) (UserResponse, error) {
	var ur UserResponse
	err := pg.db.QueryRow(context.Background(), `
	insert into users (user_name, first_name, last_name, email, pass_hash, created_at) values 
		($1,$2,$3,$4,$5,$6)
	on conflict (user_name) do 
		update set 
			user_name = coalesce(users.user_name, excluded.user_name),
			first_name = coalesce(users.first_name, excluded.first_name),
			last_name = coalesce(users.last_name, excluded.last_name),
			email = coalesce(users.email, excluded.email),
			created_at = coalesce(users.created_at, excluded.created_at),
			pass_hash = users.pass_hash
	RETURNING user_name;
	`, user.UserName, user.FirstName, user.LastName, user.Email, "pass_hash", time.Now()).Scan(&ur.UserName)
	if err != nil {
		return UserResponse{}, err
	}

	return ur, nil
}

// GetUsers ...
func (pg *PGStore) GetUsers(ctx context.Context) ([]UserResponse, error) {
	rows, err := pg.db.Query(context.Background(), "select user_name, first_name, last_name, email from users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users = make([]UserResponse, 0)
	for rows.Next() {
		var ur UserResponse
		rows.Scan(&ur.UserName, &ur.FirstName, &ur.LastName, &ur.Email)
		users = append(users, ur)
	}
	return users, nil
}
