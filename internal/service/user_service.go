package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/model"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
)

type UserService interface {
	Register(username, password string) (*model.User, error)
	Login(username, password string) (string, error) // returns JWT token
	ParseToken(tokenStr string) (uint, error)
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{repo: repo, cfg: cfg}
}

func (s *userService) Register(username, password string) (*model.User, error) {
	// check existing
	exist, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, errors.New("username already used")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Username: username,
		Password: string(hashed),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Login(username, password string) (string, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// create jwt token
	claims := jwt.RegisteredClaims{
		Subject:   string(rune(u.ID)), // we'll use ID stored in `Subject` but convert more robustly below
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.TokenTTL) * time.Second)),
	}
	// Instead of using Subject as rune conversion hack, create a map claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"exp":     claims.ExpiresAt.Unix(),
		"iat":     claims.IssuedAt.Unix(),
		"sub":     u.Username,
	})
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *userService) ParseToken(tokenStr string) (uint, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// ensure signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}
	if !tok.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	// extract user_id
	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id missing in token")
	}
	return uint(uidFloat), nil
}
