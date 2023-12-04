package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.GET ("/getwhere/:api",     /* ... */)
    router.POST("/add/:api/:json",    /* ... */)
    router.GET ("/remove/:api/:json", /* ... */)

    router.Run("localhost:8080")
}
