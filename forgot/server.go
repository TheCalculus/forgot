package forgot

import (
	"net/http"
    "github.com/gin-gonic/gin"
//  "github.com/gin-gonic/autotls"
)

type UserRegistrationInput struct {
    Name        string  `json:"name"     binding:"required"`
    Email       string  `json:"email"    binding:"required"`
    Password    string  `json:"password" binding:"required"`
}

func dummy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"response":"success"})
}

func createUser(c *gin.Context) {
    var input UserRegistrationInput

    if err := c.BindJSON(&input); err != nil {
        respondWithError(c, http.StatusBadRequest, err)
        return
    }

    user, err := CreateUser(input)
    if err != nil {
        respondWithError(c, http.StatusBadRequest, err)
        return
    }

    c.JSON(http.StatusOK, user)
}

func respondWithError(c *gin.Context, statusCode int, err error) {
    c.JSON(statusCode, gin.H{"error": err.Error()})
}

func BeginServer() {
    router := gin.Default()
    
    router.POST("/auth", createUser)
    router.GET("/auth", dummy)

    router.Run("localhost:8080")
}
