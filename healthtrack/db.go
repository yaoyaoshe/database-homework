package main

import (
    "fmt"
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    // 写死 MySQL 连接
    dsn := "root:123456@tcp(127.0.0.1:3306)/HealthTrackDB?charset=utf8mb4&parseTime=True&loc=Local"

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("连接数据库失败: %v", err)
    }

    fmt.Println("数据库连接成功！")
}
