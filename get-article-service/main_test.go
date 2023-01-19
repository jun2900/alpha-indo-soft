package main

import (
	"article-service/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockArticleRepository struct {
	mock.Mock
}

func (m *mockArticleRepository) GetAllArticle(query, author string) (results []model.Article, err error) {
	args := m.Called()
	return args.Get(0).([]model.Article), args.Error(1)
}

func TestGetArticles(t *testing.T) {
	author1 := "John"
	author2 := "Doe"
	title1 := "My first article"
	title2 := "My second article"
	body1 := "This is the body of my first article"
	body2 := "This is the body of my second article"
	now := time.Now()

	mockArticleRepo := new(mockArticleRepository)
	mockArticleRepo.On("GetAllArticle").Return([]model.Article{{
		ID:      1,
		Author:  &author1,
		Title:   &title1,
		Body:    &body1,
		Created: &now,
	}, {
		ID:      2,
		Author:  &author2,
		Title:   &title2,
		Body:    &body2,
		Created: &now,
	}}, nil)

	getArticleController := NewGetArticleController(mockArticleRepo)

	// Initialize the test server
	app := fiber.New()
	app.Get("/articles", getArticleController.handleGetArticles)

	// Test the endpoint without query parameters
	req := httptest.NewRequest("GET", "/articles", nil)
	resp, err := app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var articles []model.Article
	json.NewDecoder(resp.Body).Decode(&articles)
	assert.NotEmpty(t, articles)

	// Test the endpoint with a query parameter
	req = httptest.NewRequest("GET", "/articles?query=test", nil)
	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var filteredArticles []model.Article
	json.NewDecoder(resp.Body).Decode(&filteredArticles)
	assert.Empty(t, filteredArticles)
	for _, article := range filteredArticles {
		assert.True(t, strings.Contains(*article.Title, "test") || strings.Contains(*article.Body, "test"))
	}

	// Test the endpoint with an author parameter
	req = httptest.NewRequest("GET", "/articles?author=John", nil)
	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var filteredArticles2 []model.Article
	json.NewDecoder(resp.Body).Decode(&filteredArticles2)
	assert.NotEmpty(t, filteredArticles2)
	found := false
	for _, article := range filteredArticles2 {
		if strings.Contains(strings.ToLower(*article.Author), strings.ToLower("John")) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

//func TestGetArticlesWithMockDB(t *testing.T) {
//	// Create a mock database
//	db := new(mockDB)
//	articles := []model.Article{
//		{ID: 1, Author: "John Doe", Title: "My first article", Body: "This is the body of my first article", CreatedAt: time.Now()},
//		{ID: 2, Author: "Jane Smith", Title: "My second article", Body: "This is the body of my second article", CreatedAt: time.Now()},
//	}
//	db.On("Find", &[]Article{}, []interface{}{}).Return(db).Once().Run(func(args mock.Arguments) {
//		out := args.Get(0).(*[]Article)
//		*out = articles
//	})
//
//	// Create a request to pass to our handler
//	req, err := http.NewRequest("GET", "/articles", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//	rr := httptest.NewRecorder()
//	handler := app.Handler()
//
//	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//	// directly and pass in our Request and ResponseRecorder.
//	handler.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	// Check the response body is what we expect.
//	var responseArticles []Article
//	json.Unmarshal(rr.Body.Bytes(), &responseArticles)
//	if len(responseArticles) != 2 {
//		t.Errorf("handler returned unexpected number of articles: got %v want 2",
//			len(responseArticles))
//	}
//
//	// Ensure that the mock database's Find method was called once
//	db.AssertExpectations(t)
//}
