package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
)

func main() {
    InitDB()
    r := gin.Default()
    RegisterRoutes(r)

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    log.Printf("server listening on :%s", port)
    r.Run(":" + port)
}
