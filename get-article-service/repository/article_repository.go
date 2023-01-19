package repository

import (
	"article-service/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAllArticle(query, author string) (results []model.Article, err error)
}

func NewArticleRepository(mysqlConnection *gorm.DB) ArticleRepository {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetAllArticle(query, author string) (results []model.Article, err error) {
	db := r.mysql.Order("created desc")

	if query != "" {
		db = db.Where("title LIKE ? OR body LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if author != "" {
		db = db.Where("author = ?", author)
	}

	if err := db.Find(&results).Error; err != nil {
		return nil, ErrNotFound
	}
	return results, nil
}
