package main

import (
	"article-service/entity"
	"article-service/infrastructure"
	"article-service/model"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func main() {
	username := "john"
	password := "test123"
	dbName := "AIS"
	host := "db"
	port := "3306"

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)

	db = infrastructure.InitDb(databaseURL)
	db.AutoMigrate(&model.Article{})

	app := fiber.New()

	// Handle [POST] /articles
	app.Post("/articles", func(c *fiber.Ctx) error {
		// Parse the JSON request
		var newArticle entity.CreateArticleReq
		if err := c.BodyParser(&newArticle); err != nil {
			return c.Status(400).SendString("Invalid request")
		}

		now := time.Now()

		articleInput := model.Article{
			Author:  &newArticle.Author,
			Title:   &newArticle.Title,
			Body:    &newArticle.Body,
			Created: &now,
		}

		// Insert the new article into the MySQL database
		db.Create(&articleInput)

		// Send a success response
		return c.Status(201).JSON(articleInput)
	})

	app.Listen(":7122")
}

func createArticle(c *fiber.Ctx) error {
	// Parse the JSON request
	var newArticle entity.CreateArticleReq
	if err := c.BodyParser(&newArticle); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	now := time.Now()

	articleInput := model.Article{
		Author:  &newArticle.Author,
		Title:   &newArticle.Title,
		Body:    &newArticle.Body,
		Created: &now,
	}

	// Insert the new article into the MySQL database
	db.Create(&articleInput)

	// Send a success response
	return c.Status(201).JSON(articleInput)
}
