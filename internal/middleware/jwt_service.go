package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtService struct {
	secretKey string
	issure    string
}

func NewJWTService(secretKey string) *JwtService {
	return &JwtService{
		secretKey: secretKey,
		issure:    "marketplace-api",
	}
}

type Claim struct {
	Sum  uint     `json:"sum"`
	Role string 	`json:"role"`
	jwt.RegisteredClaims
}

func (s *JwtService) GenateToken(id uint, role string) (string, error) {
	claims := &Claim{
		Sum:  id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issure,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JwtService) ValidateToken(tokenString string) (*Claim, bool) {
	claims := &Claim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		fmt.Printf("Token validation error: %v\n", err)
		return nil, false
	}

	if !token.Valid {
		fmt.Println("Token is invalid.")
		return nil, false
	}

	return claims, true
}