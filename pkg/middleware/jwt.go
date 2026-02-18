package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

type Claims struct {
	jwt.RegisteredClaims
	Token        string         `json:"token"`
	UserId       int64          `json:"employee_id"`
	Login        string         `json:"login"`
	Role         string         `json:"role"`
	Name         string         `json:"name"`
	CustomClaims map[string]any `json:"custom_claims"`
}

const claimsKey contextKey = "claims"

func Configure(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			unauthorized(c, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			unauthorized(c, "invalid authorization header format")
			return
		}

		tokenStr := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			unauthorized(c, err.Error())
			return
		}

		// сохраняем claims в контекст
		claims.Token = tokenStr
		c.Set(string(claimsKey), claims)

		c.Next()
	}
}

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func unauthorized(c *gin.Context, detail string) {

	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse{
		Message: detail,
		Code:    "unauthorized",
	})
}

func GetClaims(c *gin.Context) *Claims {
	v, _ := c.Get(string(claimsKey))
	cl, _ := v.(*Claims)
	return cl
}
