package repository

import (
	"database/sql"
	"fmt"

	"realestate/internal/domain"
)

type FavoriteRepository struct {
    db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
    return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(userID, listingID int) error {
    query := `INSERT INTO favorites (user_id, listing_id) VALUES ($1, $2)`
    _, err := r.db.Exec(query, userID, listingID)
    return err
}

func (r *FavoriteRepository) Remove(userID, listingID int) error {
    query := `DELETE FROM favorites WHERE user_id=$1 AND listing_id=$2`
    result, err := r.db.Exec(query, userID, listingID)
    if err != nil {
        return err
    }
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("не найдено в избранном")
    }
    return nil
}

func (r *FavoriteRepository) ListByUser(userID int) ([]domain.Listing, error) {
    query := `
        SELECT
    l.id,
    l.user_id,
    l.title,
    l.price,
    l.area,
    l.rooms,
    l.city,
    l.deal_type,
    l.property_type,
    l.status,
    l.views_count,
    l.created_at
FROM favorites f
JOIN listings l ON l.id = f.listing_id
WHERE f.user_id = $1
AND l.status = 'active'
ORDER BY f.created_at DESC`

    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var listings []domain.Listing
    for rows.Next() {
    var l domain.Listing

    err := rows.Scan(
        &l.ID,
        &l.UserID,
        &l.Title,
        &l.Price,
        &l.Area,
        &l.Rooms,
        &l.City,
        &l.DealType,
        &l.PropertyType,
        &l.Status,
        &l.ViewsCount,
        &l.CreatedAt,
    )
    if err != nil {
        return nil, err
    }

    listings = append(listings, l)
}

if err := rows.Err(); err != nil {
    return nil, err
}

return listings, nil
}