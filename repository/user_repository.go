package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	if user.Role == "" {
		user.Role = "user"
	}
	_, err := ur.db.ExecContext(c,
		`INSERT INTO users (id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
	)

	return err
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	rows, err := ur.db.QueryContext(c, `SELECT id, name, email, role FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	err := ur.db.QueryRowContext(c, `SELECT id, name, email, password, role FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return user, err
	}

	return user, err
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	var user domain.User
	err := ur.db.QueryRowContext(c, `SELECT id, name, email, password, role FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return user, err
	}

	return user, err
}

func (ur *userRepository) Update(c context.Context, id string, user *domain.User) error {
	_, err := ur.db.ExecContext(c,
		`UPDATE users SET name = $1, role = $2, updated_at = NOW() WHERE id = $3`,
		user.Name,
		user.Role,
		id,
	)

	return err
}