// package main

// import (
//     "log"
//     "os"
//     "github.com/gin-gonic/gin"
// )

// func main() {
//     InitDB()
//     r := gin.Default()
//     RegisterRoutes(r)

//     port := os.Getenv("PORT")
//     if port == "" { port = "8080" }
//     log.Printf("server listening on :%s", port)
//     r.Run(":" + port)
// }
package main

import (
    "log"
    "net/http" // 必须引入这个包
    "os"
    "github.com/gin-gonic/gin"
)

// CORSMiddleware 处理跨域请求的核心配置
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        // 拦截 OPTIONS 请求，直接返回 204 No Content，不再向下传递
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

func main() {
    InitDB()
    r := gin.Default()

    // 关键步骤：启用 CORS 中间件！
    // 这行代码必须在 RegisterRoutes 之前
    r.Use(CORSMiddleware())

    RegisterRoutes(r)

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    log.Printf("server listening on :%s", port)
    r.Run(":" + port)
}