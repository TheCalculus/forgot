package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type User struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    Email     string  `json:"email"`
    Password  string  `json:"password"`
}

type UserRegistrationInput struct {
    Name      string  `json:"name"`
    Email     string  `json:"email"`
    Password  string  `json:"password"`
}

func createUser(c *gin.Context) {
    var input User

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }

    _ = User {
        Name:     input.Name,
        Password: input.Password,
    }
}

func main() {
    router := gin.Default()
    
    router.POST("/auth", createUser)

    router.Run("localhost:8080")
}
