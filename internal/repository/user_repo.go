package repository

import (
	"database/sql"
	"fmt"

	"realestate/internal/domain"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *domain.User) (*domain.User, error) {
    query := `
        INSERT INTO users (email, password, name)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`
    err := r.db.QueryRow(query, u.Email, u.Password, u.Name).
        Scan(&u.ID, &u.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    return u, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
    u := &domain.User{}
    query := `SELECT
    id,
    email,
    password,
    name,
    COALESCE(phone, ''),
    COALESCE(avatar_url, ''),
    role,
    is_verified,
    created_at
FROM users
WHERE email = $1`
    err := r.db.QueryRow(query, email).
        Scan(
    &u.ID,
    &u.Email,
    &u.Password,
    &u.Name,
    &u.Phone,
    &u.AvatarURL,
    &u.Role,
    &u.IsVerified,
    &u.CreatedAt,
)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("пользователь не найден")
    }
    return u, err
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    err := r.db.QueryRow(query, email).Scan(&exists)
    return exists, err
}
func (r *UserRepository) CreateWithRole(u *domain.User) (*domain.User, error) {
    query := `
        INSERT INTO users (email, password, name, phone, role)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, is_verified, created_at`
        err := r.db.QueryRow(query, 
        u.Email, u.Password, u.Name, u.Phone, u.Role,
    ).Scan(&u.ID, &u.IsVerified, &u.CreatedAt)
    return u, err
}
func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	u := &domain.User{}

	query := `
	SELECT
		id,
		email,
		password,
		name,
		COALESCE(phone, ''),
		COALESCE(avatar_url, ''),
		role,
		is_verified,
		created_at
	FROM users
	WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Name,
		&u.Phone,
		&u.AvatarURL,
		&u.Role,
		&u.IsVerified,
		&u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("пользователь не найден")
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}