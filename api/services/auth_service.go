package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
	"github.com/RealLiuSha/echo-admin/models/dto"
)

type options struct {
	issuer        string
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

type AuthService struct {
	opts  *options
	redis lib.Redis
}

func NewAuthService(redis lib.Redis, config lib.Config) AuthService {
	issuer := config.Name
	signingKey := fmt.Sprintf("Jwt:%s", issuer)

	opts := &options{
		issuer:        issuer,
		tokenType:     "Bearer",
		expired:       config.Auth.TokenExpired,
		signingMethod: jwt.SigningMethodHS512,
		signingKey:    []byte(signingKey),
		keyfunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.AuthTokenInvalid
			}
			return []byte(signingKey), nil
		},
	}

	return AuthService{redis: redis, opts: opts}
}

func wrapperAuthKey(key string) string {
	return fmt.Sprintf("auth:%s", key)
}

func (a AuthService) GenerateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &dto.JwtClaims{
		ID:       user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Duration(a.opts.expired) * time.Second).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	token := jwt.NewWithClaims(a.opts.signingMethod, claims)
	expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())

	err := a.redis.Set(wrapperAuthKey(claims.Username), 1, expired)
	if err != nil {
		return "", err
	}

	return token.SignedString(a.opts.signingKey)
}

func (a AuthService) ParseToken(tokenString string) (*dto.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.JwtClaims{}, a.opts.keyfunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.AuthTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.AuthTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.AuthTokenNotValidYet
			} else {
				return nil, errors.AuthTokenInvalid
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*dto.JwtClaims); ok && token.Valid {
			return claims, nil
		}
	}

	return nil, errors.AuthTokenInvalid
}

func (a AuthService) DestroyToken(username string) error {
	_, err := a.redis.Delete(wrapperAuthKey(username))
	return err
}
