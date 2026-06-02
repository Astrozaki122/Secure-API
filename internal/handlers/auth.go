package handlers

import (
	"context"
	"net/http"
	"time"

	"secure-api/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secretkey")

// -------------------- INPUT STRUCT --------------------
type AuthInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// -------------------- REGISTER --------------------
func Register(c *gin.Context) {
	var body AuthInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hash failed"})
		return
	}

	_, err = db.DB.Exec(
		context.Background(),
		"INSERT INTO users (username, password, coins) VALUES ($1, $2, 0)",
		body.Username,
		string(hash),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user already exists or db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

// -------------------- LOGIN --------------------
func Login(c *gin.Context) {
	var body AuthInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var dbPassword string

	err := db.DB.QueryRow(
		context.Background(),
		"SELECT password FROM users WHERE username=$1",
		body.Username,
	).Scan(&dbPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(body.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": body.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

// -------------------- GET COINS --------------------
func GetCoins(c *gin.Context) {
	username, exists := c.Get("username")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u := username.(string)

	var coins int

	err := db.DB.QueryRow(
		context.Background(),
		"SELECT coins FROM users WHERE username=$1",
		u,
	).Scan(&coins)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": u,
		"coins":    coins,
	})
}
