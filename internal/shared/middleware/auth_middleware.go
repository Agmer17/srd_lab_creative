package middleware

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
	RoleAdmin string     = "ADMIN"
	RoleUser  string     = "USER"
)

func getAccessToken(header string) (string, error) {

	if header == "" {
		fmt.Println("\n\n\n\n\n kosong headernya : ", header)
		return "", errors.New("Harap login terlebih dahulu sebelum mengakses fitur ini")
	}
	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Harap login terlebih dahulu sebelum mengakses fitur ini")
	}

	return parts[1], nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, err := getAccessToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(401, shared.NewErrorResponse(401, "you need to login to access this feature"))
			return
		}

		claims, err := pkg.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, shared.NewErrorResponse(401, "your session is expired. You need to re-login to access this feature!"))
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Next()
	}
}

func AuthMiddlewareFromCookie() gin.HandlerFunc {

	return func(c *gin.Context) {
		cookie, err := c.Cookie("refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(401, shared.NewErrorResponse(401, "Your session is expired. Please re-login before accessing this feature"))
			return
		}

		claims, err := pkg.VerifyRefreshToken(cookie)
		if err != nil {
			c.AbortWithStatusJSON(401, shared.NewErrorResponse(401, "Your session is expired. Please re-login before accessing this feature"))
			return
		}

		c.Set(UserIDKey, claims.UserId)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetRole(c)
		if !ok {
			c.AbortWithStatusJSON(403, "you doesn't have access to this feature")
			return
		}

		isAllowed := slices.Contains(allowedRoles, role)
		if !isAllowed {
			c.AbortWithStatusJSON(403, "you doesn't have access to this feature")
			return
		}

	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	id, ok := c.Get(UserIDKey)
	return id.(uuid.UUID), ok
}

func GetRole(c *gin.Context) (string, bool) {
	role, ok := c.Get(RoleKey)
	return role.(string), ok
}
