package middleware

import (
	
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header "})
			c.Abort()
			return
		} 
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		jwtSecret := []byte(os.Getenv("JWT_SECRET "))

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func (token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(CustomClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
		
	}
}

func generateToken(userID uint, role string) string {
	claims := CustomClaims{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return signedToken
}