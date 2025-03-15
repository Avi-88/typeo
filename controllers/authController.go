package controllers

import (
	"net/http"
	"time"
	"typeo/config"
	"typeo/models"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context) {
	var userData models.AuthInfo

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request body"})
		return
	}
	var user models.User
	config.DB.Where("username= ?",userData.UserName).Find(&user)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error":"No user account with the following username found"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid credentials"})
		return
	}
	
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256,  jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There was error generating token"})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func HandleRegistration(c *gin.Context) {
	var newUser models.NewUser

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request body"})
		return
	}
	var userFound models.User
	config.DB.Where("username= ?",newUser.Username).Find(&userFound)
	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error":"User already exists"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to register user"})
		return
	}

	user := models.User{
		Username: newUser.Username,
		Password: string(hashedPassword),
		Email: newUser.Email,
	}

	if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	c.JSON(http.StatusOK, gin.H{"user": newUser})

}