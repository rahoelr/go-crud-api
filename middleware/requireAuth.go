package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-crud/initializers"
	"go-crud/models"
	"net/http"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) {
	// Mendapatkan nilai token dari cookie pada request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
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
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Memeriksa waktu kedaluwarsa token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	expirationTime := int64(claims["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Mengambil informasi user dari database berdasarkan ID yang ada di token
	userID := uint(claims["sub"].(float64)) // Mengonversi ke tipe data uint
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Menetapkan informasi user ke dalam konteks Gin
	c.Set("user", user)

	// Lanjutkan ke handler berikutnya
	c.Next()
}
