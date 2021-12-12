package article

import (
	"time"

	"github.com/sangianpatrick/devoria-article-service/domain/account"
)

type CreateArticleResponses struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

type UpdateArticleResponses struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

type ArticleResponses struct {
	ID             int64
	Title          string
	Subtitle       string
	Content        string
	Status         ArticleStatus
	CreatedAt      time.Time
	PublishedAt    *time.Time
	LastModifiedAt *time.Time
	Author         account.Account
}
