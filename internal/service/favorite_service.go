package service

import (
    "realestate/internal/domain"
    "realestate/internal/repository"
)

type FavoriteService struct {
    repo *repository.FavoriteRepository
}

func NewFavoriteService(repo *repository.FavoriteRepository) *FavoriteService {
    return &FavoriteService{repo: repo}
}

func (s *FavoriteService) Add(userID, listingID int) error {
    return s.repo.Add(userID, listingID)
}

func (s *FavoriteService) Remove(userID, listingID int) error {
    return s.repo.Remove(userID, listingID)
}

func (s *FavoriteService) List(userID int) ([]domain.Listing, error) {
    return s.repo.ListByUser(userID)
}