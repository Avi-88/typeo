package controllers

import (
	// "fmt"
	"log"
	"net/http"
	"typeo/config"
	"typeo/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context)  {
	var users []models.User

	err := config.DB.Select("username", "email").Find(&users).Error
	if err != nil {
		log.Fatal("Error fetching users")
		c.JSON(500,"There was some error fetching users")
	}
	c.JSON(http.StatusOK, users)
}

// func createUser(c *gin.Context) {
// 	var newUser models.NewUser

//     if err := c.ShouldBindJSON(&newUser); err != nil {
// 		fmt.Println("Invalid request body received - user")
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
// 	var userExists models.User
//     config.DB.Where("userName=?", newUser.Username).Find(&userExists)
// 	if userExists.ID != 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "This username is already in use"})
// 		return
// 	}

// 	if err := config.DB.Create(&newUser).Error; err != nil {
// 		fmt.Println("There was an error creating the user")
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
// 	}
// 	c.JSON(http.StatusOK, &newUser)
// }


func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	var user models.User

    if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This user doesnt exist"})
		return
	}

	var updatedUser map[string]interface{}
	 if err:= c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invald request body format"})
		return
	 }

	if err := config.DB.Model(&user).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"There was an error in updating the userdata"})
		return
	}

	c.JSON(http.StatusOK, user)
}


func DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	var user models.User

    if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This user doesnt exist"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"There was an error in deleting the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message":"User was deleted successfully!"})
}