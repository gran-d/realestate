package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"realestate/internal/domain"
	"realestate/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtSecret string
}

func NewUserService(
	repo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
	secret string,
) *UserService {
	return &UserService{
		repo:      repo,
		tokenRepo: tokenRepo,
		jwtSecret: secret,
	}
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

func (s *UserService) Register(email, password, name string) (*domain.User, error) {

	exists, _ := s.repo.ExistsByEmail(email)
	if exists {
		return nil, fmt.Errorf("email уже используется")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hash),
		Name:     name,
	}

	return s.repo.Create(user)
}

func (s *UserService) RegisterWithRole(email, password, name, phone, role string) (*domain.User, error) {

	exists, _ := s.repo.ExistsByEmail(email)
	if exists {
		return nil, fmt.Errorf("email уже используется")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hash),
		Name:     name,
		Phone:    phone,
		Role:     role,
	}

	return s.repo.CreateWithRole(user)
}

func (s *UserService) Login(email, password string) (*LoginResponse, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = s.tokenRepo.Save(user.ID, refreshToken, expiresAt)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *UserService) generateAccessToken(userID int, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func generateRandomToken() (string, error) {

	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (s *UserService) Refresh(refreshToken string) (string, error) {

	rt, err := s.tokenRepo.Find(refreshToken)
	if err != nil {
		return "", fmt.Errorf("невалидный или истёкший refresh token")
	}
    
	user, err := s.repo.GetByID(rt.UserID)
	if err != nil {
		return "", fmt.Errorf("пользователь не найден")
	}

	return s.generateAccessToken(user.ID, user.Role)
}

func (s *UserService) Logout(refreshToken string) error {
	return s.tokenRepo.Delete(refreshToken)
}

func (s *UserService) LogoutAll(userID int) error {
	return s.tokenRepo.DeleteAllForUser(userID)
}