package main

import (
	"article-service/entity"
	"article-service/infrastructure"
	"article-service/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPostArticles(t *testing.T) {
	username := "root"
	password := "test123"
	dbName := "AIS"
	host := "localhost"
	port := "3055"

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)

	db = infrastructure.InitDb(databaseURL)
	db.AutoMigrate(&model.Article{})

	// Set up the test server
	app := fiber.New()
	app.Post("/articles", createArticle) // reuse the handler function from main

	// Prepare the sample data
	sampleData := entity.CreateArticleReq{
		Author: "John Doe",
		Title:  "Test Title",
		Body:   "Test Body",
	}
	jsonData, _ := json.Marshal(sampleData)

	// Send the POST request
	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the response status code
	if res.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	// Assert the response body
	var article model.Article
	fmt.Println(article)
	json.NewDecoder(res.Body).Decode(&article)
	if *article.Author != sampleData.Author || *article.Title != sampleData.Title || *article.Body != sampleData.Body {
		t.Errorf("Unexpected response body: %+v", article)
	}
}
