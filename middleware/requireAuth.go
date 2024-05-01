package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-crud/initializers"
	"go-crud/models"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) error {
	// Mendapatkan nilai token dari cookie pada request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		return fmt.Errorf("missing Authorization cookie")
	}

	// Parse token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Memastikan bahwa metode penandatanganan yang digunakan adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Mengembalikan kunci untuk verifikasi token
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("invalid or expired token")
	}

	// Memeriksa waktu kedaluwarsa token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}
	expirationTime := int64(claims["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		return fmt.Errorf("expired token")
	}

	// Mengambil informasi user dari database berdasarkan ID yang ada di token
	userID := uint(claims["sub"].(float64)) // Mengonversi ke tipe data uint
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Menetapkan informasi user ke dalam konteks Gin
	c.Set("user", user)

	return nil
}
