package service

import (
	"fmt"

	"realestate/internal/domain"
	"realestate/internal/repository"
)

type ListingService struct {
    repo *repository.ListingRepository
}

func NewListingService(repo *repository.ListingRepository) *ListingService {
    return &ListingService{repo: repo}
}

func (s *ListingService) Create(l *domain.Listing) (*domain.Listing, error) {
    if l.Title == "" {
        return nil, fmt.Errorf("заголовок обязателен")
    }
    if l.Price <= 0 {
        return nil, fmt.Errorf("цена должна быть больше нуля")
    }
    if l.Area <= 0 {
        return nil, fmt.Errorf("площадь должна быть больше нуля")
    }
    if l.City == "" {
        return nil, fmt.Errorf("город обязателен")
    }
    if l.DealType != "sale" && l.DealType != "rent" {
        return nil, fmt.Errorf("тип должен быть 'sale' или 'rent'")
    }
    validProperties := map[string]bool{
        "apartment": true,
        "house":     true,
        "commercial": true,
    }
    if !validProperties[l.PropertyType] {
        return nil, fmt.Errorf("property: apartment, house или commercial")
    }
    return s.repo.Create(l)
}

func (s *ListingService) GetByID(id int) (*domain.Listing, error) {
    return s.repo.GetByID(id)
}

func (s *ListingService) Search(f domain.ListingFilter) ([]domain.Listing, error) {
    return s.repo.Search(f)
}

func (s *ListingService) Update(l *domain.Listing, requestingUserID int) error {
    existing, err := s.repo.GetByIDNoView(l.ID)
    if err != nil {
        return fmt.Errorf("объявление не найдено")
    }
    if existing.UserID != requestingUserID {
        return fmt.Errorf("нет прав на редактирование")
    }
    l.UserID = requestingUserID
    return s.repo.Update(l)
}

func (s *ListingService) Delete(id, userID int) error {
    return s.repo.Delete(id, userID)
}