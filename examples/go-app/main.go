package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
)

type User struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

type Post struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID uint   `json:"user_id"`
}

var db *gorm.DB

func init() {
	godotenv.Load()
	
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	db.AutoMigrate(&User{}, &Post{})
}

func main() {
	defer db.Close()
	
	router := gin.Default()
	
	router.GET("/api/users", getUsers)
	router.POST("/api/users", createUser)
	router.GET("/api/users/:id", getUser)
	router.GET("/api/posts", getPosts)
	router.POST("/api/posts", createPost)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server running on port %s\n", port)
	router.Run(":" + port)
}

func getUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(200, users)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	db.Create(&user)
	c.JSON(201, user)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(200, user)
}

func getPosts(c *gin.Context) {
	var posts []Post
	db.Find(&posts)
	c.JSON(200, posts)
}

func createPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	db.Create(&post)
	c.JSON(201, post)
}
