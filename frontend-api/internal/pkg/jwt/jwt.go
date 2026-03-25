package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWT struct {
	secret     string
	expiration int64
	issuer     string
	redis      *redis.Client
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func New(redisClient *redis.Client) *JWT {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	v.ReadInConfig()

	return &JWT{
		secret:     v.GetString("jwt.secret"),
		expiration: v.GetInt64("jwt.expiration"),
		issuer:     v.GetString("jwt.issuer"),
		redis:      redisClient,
	}
}

func NewWithConfig(secret string, expiration int64, issuer string, redisClient *redis.Client) *JWT {
	return &JWT{
		secret:     secret,
		expiration: expiration,
		issuer:     issuer,
		redis:      redisClient,
	}
}

// GenerateToken generates a new JWT token
func (j *JWT) GenerateToken(userID int64, login, name, avatarURL, email string) (string, string, error) {
	now := time.Now()
	jti := fmt.Sprintf("%d-%d", userID, now.UnixNano())

	claims := Claims{
		UserID:    userID,
		Login:     login,
		Name:      name,
		AvatarURL: avatarURL,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.expiration) * time.Second)),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}

	return signedToken, jti, nil
}

// ValidateToken validates a JWT token and checks if it's blacklisted
func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if token is blacklisted
	if j.isBlacklisted(claims.ID) {
		return nil, fmt.Errorf("token has been revoked")
	}

	return claims, nil
}

// AddToBlacklist adds a token's jti to the Redis blacklist
func (j *JWT) AddToBlacklist(ctx context.Context, jti string, exp time.Time) error {
	if j.redis == nil {
		return nil
	}

	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}

	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	return j.redis.Set(ctx, key, "1", ttl).Err()
}

// isBlacklisted checks if a token's jti is in the Redis blacklist
func (j *JWT) isBlacklisted(jti string) bool {
	if j.redis == nil {
		return false
	}

	ctx := context.Background()
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	result, err := j.redis.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return result > 0
}

// GetExpiration returns the token expiration time in seconds
func (j *JWT) GetExpiration() int64 {
	return j.expiration
}
