package libecho

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/helloferdie/golib/libjwt"
	"github.com/labstack/echo/v4"
)

// GetRealIP - Get Real IP in reverse proxy environment
func GetRealIP(c echo.Context) string {
	if c != nil {
		ip := c.Request().Header.Get("X-Real-Ip")
		if ip == "" {
			ip = c.Request().RemoteAddr
		}
		return ip
	}
	return ""
}

// GetJWTClaims - Get JWT claims from middleware
func GetJWTClaims(c echo.Context) jwt.MapClaims {
	test := c.Get("user")
	if test == nil {
		return GetJWTClaimsHeader(c)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}

// GetJWTClaimsHeader - Get JWT claims from header
func GetJWTClaimsHeader(c echo.Context) jwt.MapClaims {
	header := c.Request().Header.Get("Authorization")
	if header != "" {
		auth := strings.Split(header, "Bearer ")
		if len(auth) < 2 {
			return nil
		}

		tokenValue := auth[1]
		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			return libjwt.GetByte(), nil
		})
		if err == nil && token.Valid {
			return token.Claims.(jwt.MapClaims)
		}
	}
	return nil
}
