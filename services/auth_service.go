package services

import (
	"errors"
	"time"

	"github.com/galiherlangga/clinic-appointment/configs"
	"github.com/galiherlangga/clinic-appointment/models"
	"github.com/galiherlangga/clinic-appointment/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Login(email, password string) (*TokenPair, error)
	Register(email, password, fullName string) error
	GetMe(userID uint) (*models.User, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	Logout(token string) error
}

type authService struct {
	userRepo      repositories.UserRepository
	blacklistRepo BlacklistService
}

func NewAuthService(userRepo repositories.UserRepository, blacklistRepo BlacklistService) AuthService {
	return &authService{userRepo: userRepo, blacklistRepo: blacklistRepo}
}

func (s *authService) Login(email, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenPair(user)
}

func (s *authService) generateTokenPair(user *models.User) (*TokenPair, error) {
	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(configs.AppConfig.JWTExpirationHours)).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(configs.AppConfig.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * time.Duration(configs.AppConfig.JWTRefreshExpirationHours)).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(configs.AppConfig.JWTRefreshSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *authService) Register(email, password, fullName string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
		Role:     models.RoleUser,
	}

	return s.userRepo.Create(user)
}

func (s *authService) GetMe(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.AppConfig.JWTRefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := uint(claims["user_id"].(float64))
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateTokenPair(user)
}

func (s *authService) Logout(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	exp := int64(claims["exp"].(float64))
	expiration := time.Unix(exp, 0).Sub(time.Now())

	if expiration > 0 {
		return s.blacklistRepo.BlacklistToken(tokenString, expiration)
	}

	return nil
}
