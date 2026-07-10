package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"realestate/internal/domain"
)

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) Create(l *domain.Listing) (*domain.Listing, error) {
	query := `
		INSERT INTO listings
			(user_id, title, description, price, area, rooms, floor, total_floors,
			 city, address, latitude, longitude, deal_type, property_type)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, status, views_count, created_at, updated_at`
	err := r.db.QueryRow(query,
		l.UserID, l.Title, l.Description, l.Price, l.Area, l.Rooms,
		l.Floor, l.TotalFloors, l.City, l.Address,
		l.Latitude, l.Longitude, l.DealType, l.PropertyType,
	).Scan(&l.ID, &l.Status, &l.ViewsCount, &l.CreatedAt, &l.UpdatedAt)
	return l, err
}

func (r *ListingRepository) GetByID(id int) (*domain.Listing, error) {
	query := `
		WITH updated AS (
			UPDATE listings SET views_count = views_count + 1
			WHERE id = $1 AND status != 'inactive'
			RETURNING *
		)
		SELECT id, user_id, title, description, price, area, rooms,
		       COALESCE(floor, 0), COALESCE(total_floors, 0),
		       city, COALESCE(address, ''),
		       COALESCE(latitude, 0), COALESCE(longitude, 0),
		       deal_type, property_type, status, views_count,
		       created_at, updated_at
		FROM updated`

	l := &domain.Listing{}
	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.UserID, &l.Title, &l.Description, &l.Price, &l.Area,
		&l.Rooms, &l.Floor, &l.TotalFloors, &l.City, &l.Address,
		&l.Latitude, &l.Longitude, &l.DealType, &l.PropertyType,
		&l.Status, &l.ViewsCount, &l.CreatedAt, &l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("объявление не найдено")
	}
	if err != nil {
    	return nil, fmt.Errorf("ошибка scan: %w", err)
	}
	return l, nil
}

func (r *ListingRepository) GetByIDNoView(id int) (*domain.Listing, error) {
	query := `
		SELECT id, user_id, title, description, price, area, rooms, floor, total_floors,
		       city, address, latitude, longitude, deal_type, property_type,
		       status, views_count, created_at, updated_at
		FROM listings WHERE id = $1`
	l := &domain.Listing{}
	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.UserID, &l.Title, &l.Description, &l.Price, &l.Area,
		&l.Rooms, &l.Floor, &l.TotalFloors, &l.City, &l.Address,
		&l.Latitude, &l.Longitude, &l.DealType, &l.PropertyType,
		&l.Status, &l.ViewsCount, &l.CreatedAt, &l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("объявление не найдено")
	}
	return l, err
}

func (r *ListingRepository) Search(f domain.ListingFilter) ([]domain.Listing, error) {
	where := []string{"status = 'active'"}
	args := []interface{}{}
	i := 1

	if f.City != "" {
		where = append(where, fmt.Sprintf("city ILIKE $%d", i))
		args = append(args, "%"+f.City+"%")
		i++
	}
	if f.DealType != "" {
		where = append(where, fmt.Sprintf("deal_type = $%d", i))
		args = append(args, f.DealType)
		i++
	}
	if f.PropertyType != "" {
		where = append(where, fmt.Sprintf("property_type = $%d", i))
		args = append(args, f.PropertyType)
		i++
	}
	if f.MinPrice > 0 {
		where = append(where, fmt.Sprintf("price >= $%d", i))
		args = append(args, f.MinPrice)
		i++
	}
	if f.MaxPrice > 0 {
		where = append(where, fmt.Sprintf("price <= $%d", i))
		args = append(args, f.MaxPrice)
		i++
	}
	if f.Rooms > 0 {
		where = append(where, fmt.Sprintf("rooms = $%d", i))
		args = append(args, f.Rooms)
		i++
	}
	if f.MinArea > 0 {
		where = append(where, fmt.Sprintf("area >= $%d", i))
		args = append(args, f.MinArea)
		i++
	}
	if f.MaxArea > 0 {
		where = append(where, fmt.Sprintf("area <= $%d", i))
		args = append(args, f.MaxArea)
		i++
	}

	limit := f.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := 0
	if f.Page > 1 {
		offset = (f.Page - 1) * limit
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, title, price, area, rooms, city,
		       deal_type, property_type, status, views_count, created_at
		FROM listings
		WHERE %s
		ORDER BY created_at DESC
		LIMIT %d OFFSET %d`,
		strings.Join(where, " AND "), limit, offset,
	)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []domain.Listing
	for rows.Next() {
		var l domain.Listing
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Title, &l.Price, &l.Area,
			&l.Rooms, &l.City, &l.DealType, &l.PropertyType,
			&l.Status, &l.ViewsCount, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, rows.Err()
}

func (r *ListingRepository) Update(l *domain.Listing) error {
	query := `
		UPDATE listings
		SET title=$1, description=$2, price=$3, area=$4, rooms=$5,
		    floor=$6, total_floors=$7, city=$8, address=$9,
		    latitude=$10, longitude=$11, deal_type=$12, property_type=$13,
		    updated_at=NOW()
		WHERE id=$14 AND user_id=$15`
	result, err := r.db.Exec(query,
		l.Title, l.Description, l.Price, l.Area, l.Rooms,
		l.Floor, l.TotalFloors, l.City, l.Address,
		l.Latitude, l.Longitude, l.DealType, l.PropertyType,
		l.ID, l.UserID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("объявление не найдено или нет прав")
	}
	return nil
}

func (r *ListingRepository) UpdateStatus(id, userID int, status string) error {
	result, err := r.db.Exec(
		`UPDATE listings SET status=$1, updated_at=NOW() WHERE id=$2 AND user_id=$3`,
		status, id, userID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("объявление не найдено или нет прав")
	}
	return nil
}

func (r *ListingRepository) Delete(id, userID int) error {
	return r.UpdateStatus(id, userID, domain.StatusInactive)
}

func (r *ListingRepository) GetByUserID(userID int) ([]domain.Listing, error) {
	query := `
		SELECT id, user_id, title, price, area, rooms, city,
		       deal_type, property_type, status, views_count, created_at
		FROM listings WHERE user_id = $1
		ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []domain.Listing
	for rows.Next() {
		var l domain.Listing
		rows.Scan(
			&l.ID, &l.UserID, &l.Title, &l.Price, &l.Area,
			&l.Rooms, &l.City, &l.DealType, &l.PropertyType,
			&l.Status, &l.ViewsCount, &l.CreatedAt,
		)
		listings = append(listings, l)
	}
	return listings, rows.Err()
}

func (r *ListingRepository) AddPhoto(photo *domain.ListingPhoto) (*domain.ListingPhoto, error) {
	query := `
		INSERT INTO listing_photos (listing_id, url, is_main, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	err := r.db.QueryRow(query,
		photo.ListingID, photo.URL, photo.IsMain, photo.SortOrder,
	).Scan(&photo.ID, &photo.CreatedAt)
	return photo, err
}

func (r *ListingRepository) GetPhotos(listingID int) ([]domain.ListingPhoto, error) {
	rows, err := r.db.Query(
		`SELECT id, listing_id, url, is_main, sort_order, created_at
		 FROM listing_photos WHERE listing_id = $1 ORDER BY sort_order`,
		listingID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []domain.ListingPhoto
	for rows.Next() {
		var p domain.ListingPhoto
		rows.Scan(&p.ID, &p.ListingID, &p.URL, &p.IsMain, &p.SortOrder, &p.CreatedAt)
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

func (r *ListingRepository) DeletePhoto(photoID, listingID, userID int) error {
	query := `
		DELETE FROM listing_photos
		WHERE id = $1 AND listing_id = $2
		  AND listing_id IN (SELECT id FROM listings WHERE user_id = $3)`
	result, err := r.db.Exec(query, photoID, listingID, userID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("фото не найдено или нет прав")
	}
	return nil
}

func (r *ListingRepository) AdminDelete(id int) error {
	result, err := r.db.Exec(
		`UPDATE listings SET status='inactive', updated_at=NOW() WHERE id=$1`, id,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("объявление не найдено")
	}
	return nil
}

func (r *ListingRepository) GetStats() (map[string]interface{}, error) {
	var activeListings, totalUsers, totalAgents, totalFavorites, totalViews int

	r.db.QueryRow(`SELECT COUNT(*) FROM listings WHERE status='active'`).Scan(&activeListings)
	r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE role='agent'`).Scan(&totalAgents)
	r.db.QueryRow(`SELECT COUNT(*) FROM favorites`).Scan(&totalFavorites)
	r.db.QueryRow(`SELECT COALESCE(SUM(views_count),0) FROM listings`).Scan(&totalViews)

	return map[string]interface{}{
		"active_listings":  activeListings,
		"total_users":      totalUsers,
		"total_agents":     totalAgents,
		"total_favorites":  totalFavorites,
		"total_views":      totalViews,
	}, nil
}