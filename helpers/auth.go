package helpers

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/PeteCroton/go-basic/modules/middlewares/models"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

// Middleware JWT function
func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
	})
}

type gobasicAuth struct {
	mapClaims *gobasicMapClaims
}

type gobasicAdmin struct {
	*gobasicAuth
}

type gobasicApiKey struct {
	*gobasicAuth
}

type gobasicMapClaims struct {
	Claims *models.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IGobasicAuth interface {
	SignToken() string
}

type IGobasicAdmin interface {
	SignToken() string
}

type IGobasicApiKey interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *gobasicAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY"))) //token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return ss
}

func (a *gobasicAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString([]byte(os.Getenv("APP_ADMIN_KEY"))) //token.SignedString(os.Getenv("APP_ADMIN_KEY"))
	return ss
}

func (a *gobasicApiKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString([]byte(os.Getenv("APP_API_KEY"))) //token.SignedString(os.Getenv("APP_API_KEY"))
	return ss
}

func ParseToken(tokenString string) (*gobasicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &gobasicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return os.Getenv("JWT_SECRET_KEY"), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*gobasicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(tokenString string) (*gobasicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &gobasicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return os.Getenv("APP_ADMIN_KEY"), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*gobasicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseApiKey(tokenString string) (*gobasicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &gobasicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return os.Getenv("APP_API_KEY"), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*gobasicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(claims *models.UserClaims, exp int64) string {
	obj := &gobasicAuth{
		mapClaims: &gobasicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gobasicshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewGobasicAuth(tokenType TokenType, claims *models.UserClaims) (IGobasicAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(claims), nil
	case Refresh:
		return newRefreshToken(claims), nil
	case Admin:
		return newAdminToken(), nil
	case ApiKey:
		return newApiKey(), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(claims *models.UserClaims) IGobasicAuth {
	jwt_access_expires := os.Getenv("JWT_ACCESS_EXPIRES")
	time_duration, err := strconv.Atoi(jwt_access_expires)

	if err != nil {
		fmt.Println("Error during conversion")
		return nil
	}

	return &gobasicAuth{
		mapClaims: &gobasicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gobasicshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(time_duration),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(claims *models.UserClaims) IGobasicAuth {
	jwt_access_expires := os.Getenv("JWT_ACCESS_EXPIRES")
	time_duration, err := strconv.Atoi(jwt_access_expires)

	if err != nil {
		fmt.Println("Error during conversion")
		return nil
	}
	return &gobasicAuth{
		mapClaims: &gobasicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gobasicshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(time_duration),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken() IGobasicAuth {
	return &gobasicAdmin{
		gobasicAuth: &gobasicAuth{
			mapClaims: &gobasicMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "gobasicshop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func newApiKey() IGobasicAuth {
	return &gobasicApiKey{
		gobasicAuth: &gobasicAuth{
			mapClaims: &gobasicMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    os.Getenv("APP_NAME"),
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
