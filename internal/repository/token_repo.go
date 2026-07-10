package repository

import (
    "database/sql"
    "fmt"
    "time"

    "realestate/internal/domain"
)

type TokenRepository struct {
    db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
    return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(userID int, token string, expiresAt time.Time) error {
    _, err := r.db.Exec(
        `INSERT INTO refresh_tokens (user_id, token, expires_at)
         VALUES ($1, $2, $3)`,
        userID, token, expiresAt,
    )
    return err
}

func (r *TokenRepository) Find(token string) (*domain.RefreshToken, error) {
    rt := &domain.RefreshToken{}
    err := r.db.QueryRow(
        `SELECT id, user_id, token, expires_at, created_at
         FROM refresh_tokens
         WHERE token = $1 AND expires_at > NOW()`,
        token,
    ).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("токен не найден или истёк")
    }
    return rt, err
}

func (r *TokenRepository) Delete(token string) error {
    _, err := r.db.Exec(
        `DELETE FROM refresh_tokens WHERE token = $1`, token,
    )
    return err
}

func (r *TokenRepository) DeleteAllForUser(userID int) error {
    _, err := r.db.Exec(
        `DELETE FROM refresh_tokens WHERE user_id = $1`, userID,
    )
    return err
}

func (r *TokenRepository) DeleteExpired() error {
    _, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE expires_at < NOW()`)
    return err
}