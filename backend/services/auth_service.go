package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/harih/pulse-poll/backend/models"
)

type AuthService struct {
	Users     *mongo.Collection
	JWTSecret string
}

func NewAuthService(db *mongo.Database, jwtSecret string) *AuthService {
	return &AuthService{
		Users:     db.Collection("users"),
		JWTSecret: jwtSecret,
	}
}

func (s *AuthService) Signup(ctx context.Context, name, email, password string) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	var existing models.User
	err := s.Users.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == nil {
		return nil, errors.New("email already registered")
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	result, err := s.Users.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID

	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User

	err := s.Users.FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil, errors.New("invalid email or password")
		}

		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return signedToken, &user, nil
}
