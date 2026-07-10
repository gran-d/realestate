package domain

import "time"

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone,omitempty"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type Listing struct {
	ID           int            `json:"id"`
	UserID       int            `json:"user_id"`
	Title        string         `json:"title"`
	Description  string         `json:"description,omitempty"`
	Price        int64          `json:"price"`
	Area         float64        `json:"area"`
	Rooms        int            `json:"rooms"`
	Floor        int            `json:"floor,omitempty"`
	TotalFloors  int            `json:"total_floors,omitempty"`
	City         string         `json:"city"`
	Address      string         `json:"address,omitempty"`
	Latitude     float64        `json:"latitude,omitempty"`
	Longitude    float64        `json:"longitude,omitempty"`
	DealType     string         `json:"deal_type"`
	PropertyType string         `json:"property_type"`
	Status       string         `json:"status"`
	ViewsCount   int            `json:"views_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Photos       []ListingPhoto `json:"photos,omitempty"`
}

type ListingPhoto struct {
	ID        int       `json:"id"`
	ListingID int       `json:"listing_id"`
	URL       string    `json:"url"`
	IsMain    bool      `json:"is_main"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type ListingFilter struct {
	City         string
	DealType     string
	PropertyType string
	MinPrice     int64
	MaxPrice     int64
	Rooms        int
	MinArea      float64
	MaxArea      float64
	Status       string
	Page         int
	Limit        int
}

type Favorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ListingID int       `json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID         int       `json:"id"`
	ReviewerID int       `json:"reviewer_id"`
	AgentID    int       `json:"agent_id"`
	ListingID  int       `json:"listing_id,omitempty"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	ListingID  int       `json:"listing_id,omitempty"`
	Text       string    `json:"text"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type SearchHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	RoleBuyer = "buyer"
	RoleAgent = "agent"
	RoleAdmin = "admin"
)

const (
	DealTypeSale = "sale"
	DealTypeRent = "rent"
)

const (
	PropertyApartment  = "apartment"
	PropertyHouse      = "house"
	PropertyCommercial = "commercial"
	PropertyLand       = "land"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusSold     = "sold"
	StatusRented   = "rented"
)
type RefreshToken struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}