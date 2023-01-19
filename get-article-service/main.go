package main

import (
	"article-service/get-article-service/repository"
	"article-service/infrastructure"
	"article-service/model"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	ctx = context.Background()

	// Connect to Redis
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
)

type getArticlesController struct {
	ArticleRepository repository.ArticleRepository
}

func NewGetArticleController(articleRepository repository.ArticleRepository) *getArticlesController {
	return &getArticlesController{
		ArticleRepository: articleRepository,
	}
}

func main() {
	username := "john"
	password := "test123"
	dbName := "AIS"
	host := "db"
	port := "3306"

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)

	db = infrastructure.InitDb(databaseURL)
	db.AutoMigrate(&model.Article{})

	//repository
	articleRepository := repository.NewArticleRepository(db)

	articleController := NewGetArticleController(articleRepository)

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowMethods: "GET"}))

	// Handle [GET] /articles
	app.Get("/articles", articleController.handleGetArticles)

	app.Listen(":7123")
}

func (g *getArticlesController) handleGetArticles(c *fiber.Ctx) error {
	// Check for query parameters
	query := c.Query("query")
	author := c.Query("author")

	// Try to get the articles from the Redis cache
	var articles []model.Article
	val, err := client.Get(ctx, "articles").Result()
	if err == nil {
		err := json.Unmarshal([]byte(val), &articles)
		if err != nil {
			panic(err)
		}

		//filter articles based on query
		if query != "" {
			for i, article := range articles {
				if !strings.Contains(*article.Title, query) || !strings.Contains(*article.Body, query) {
					if len(articles) == 1 {
						articles = nil
					} else {
						articles = append(articles[:i], articles[i+1:]...)
					}
				}
			}
		}

		//filter articles based on authors name
		if author != "" {
			for i, article := range articles {
				if !strings.Contains(strings.ToLower(*article.Author), strings.ToLower(author)) {
					if len(articles) == 1 {
						articles = nil
					} else {
						articles = append(articles[:i], articles[i+1:]...)
					}
				}
			}
		}

		return c.JSON(articles)
	}

	// If the articles are not in the cache, get them from the MySQL database
	articles, err = g.ArticleRepository.GetAllArticle(query, author)
	if err != nil {
		panic(err)
	}

	// Store the articles in the Redis cache
	json, err := json.Marshal(articles)
	if err != nil {
		panic(err)
	}
	client.Set(ctx, "articles", json, time.Second*5)

	// Send the articles as a JSON response
	return c.JSON(articles)
}
