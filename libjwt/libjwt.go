package libjwt

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config - JWT configuration
type Config struct {
	Secret             string
	SecretByte         []byte
	TokenExpiry        time.Duration
	RefreshTokenExpiry time.Duration
}

var initialize = false
var cfg = Config{}

// loadConfig - Load initial configuration
func loadConfig() {
	if !initialize {
		cfg.Secret = os.Getenv("jwt_secret")
		cfg.SecretByte = []byte(cfg.Secret)

		expiry, err := strconv.ParseInt(os.Getenv("jwt_expiry_minutes"), 10, 64)
		if err != nil {
			expiry = 15
		}
		cfg.TokenExpiry = time.Duration(expiry)

		refreshExpiry, err := strconv.ParseInt(os.Getenv("jwt_expiry_refresh_days"), 10, 64)
		if err != nil {
			refreshExpiry = 30
		}
		cfg.RefreshTokenExpiry = time.Duration(refreshExpiry)
		initialize = true
	}
}

// GetByte - Get JWT secret in byte format
func GetByte() []byte {
	loadConfig()
	return cfg.SecretByte
}

// Generate -
func Generate(dtMain map[string]interface{}, dtRefresh map[string]interface{}) map[string]interface{} {
	loadConfig()
	expiry := time.Now().UTC().Add(time.Minute * cfg.TokenExpiry)
	//expiry := time.Now().UTC().Add(time.Second * 10) // Shorter JWT token, for debug purpose
	refreshExpiry := time.Now().UTC().Add(time.Hour * 24 * cfg.RefreshTokenExpiry)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range dtMain {
		claims[k] = v
	}
	claims["exp"] = expiry.Unix()

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	for k, v := range dtRefresh {
		refreshClaims[k] = v
	}
	refreshClaims["exp"] = refreshExpiry.Unix()

	t, _ := token.SignedString(cfg.SecretByte)
	rt, _ := refreshToken.SignedString(cfg.SecretByte)
	return map[string]interface{}{
		"token":                t,
		"token_expiry":         claims["exp"],
		"refresh_token":        rt,
		"refresh_token_expiry": refreshClaims["exp"],
	}
}

// Model -
type Model struct {
	ID     string
	UserID int64
	Access string
}

// Parse -
func Parse(c jwt.MapClaims) *Model {
	m := new(Model)

	vInt, ok := c["user_id"].(float64)
	if ok {
		m.UserID = int64(vInt)
	}

	vStr, ok := c["access"].(string)
	if ok {
		m.Access = vStr
	}

	vStr, ok = c["id"].(string)
	if ok {
		m.ID = vStr
	}
	return m
}
