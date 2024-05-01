package controllers

import (
	"github.com/gin-gonic/gin"
	"go-crud/initializers"
	"go-crud/middleware"
	"go-crud/models"
	"net/http"
)

func PostCreate(c *gin.Context) {

	// Get Data off req body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	err := middleware.RequireAuth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User belum login",
		})
		return
	}

	// Create a post
	post := models.Post{Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}
	// return it
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostIndex(c *gin.Context) {
	// Gunakan middleware RequireAuth
	err := middleware.RequireAuth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User belum login",
		})
		return
	}

	var posts []models.Post
	initializers.DB.Find(&posts)

	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostShow(c *gin.Context) {

	err := middleware.RequireAuth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User belum login",
		})
		return
	}

	id := c.Param("id")
	var post models.Post
	initializers.DB.First(&post, id)

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {

	err := middleware.RequireAuth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User belum login",
		})
		return
	}

	id := c.Param("id")

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	var post models.Post
	initializers.DB.First(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostDelete(c *gin.Context) {
	// get the id
	id := c.Param("id")

	err := middleware.RequireAuth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User belum login",
		})
		return
	}

	// delete posts
	initializers.DB.Delete(&models.Post{}, id)

	//respone
	c.JSON(200, gin.H{
		"message": "Success Delete Post",
	})
}
