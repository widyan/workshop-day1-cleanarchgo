package unittest

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/domain/article"
	"github.com/sangianpatrick/devoria-article-service/response"
	"github.com/stretchr/testify/assert"
)

/*
func TestSave(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	cfg := config.New()
	jsonWebToken := jwt.NewJSONWebToken(jwt.GetRSAPrivateKey("../../../secret/id_rsa"), jwt.GetRSAPublicKey("../../../secret/id_rsa.pub"))
	encryption := crypto.NewAES256CBC(cfg.AES.SecretKey)
	// sess := session.NewRedisSessionStoreAdapter(*redis.Client, time.Hour*24*1)

	articleRepository := new(MockNewArticleRepository)

	ctx := context.Background()

	var resp response.Response
	request := article.CreateArticleRequest{
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}

	articleRepository.On("Save", ctx, article.Article{
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}).Return(13, nil)
	articleUsecase := article.NewArticleUsecase(cfg.GlobalIV, nil, jsonWebToken, encryption, location, articleRepository)
	resp = articleUsecase.Save(ctx, 0, request)
	log.Println(resp)

	assert.Equal(t, resp, response.Success(response.StatusCreated, request))
}
*/

func TestSave(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	articleRepository := new(MockNewArticleRepository)
	ctx := context.Background()

	var resp response.Response
	request := article.CreateArticleRequest{
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}

	articleRepository.On("Save", ctx, article.Article{
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}).Return(13, nil)
	articleUsecase := article.NewArticleUsecase(nil, location, articleRepository)
	resp = articleUsecase.Save(ctx, 0, request)
	log.Println(resp)

	assert.Equal(t, resp, response.Success(response.StatusCreated, article.CreateArticleResponses{
		ID:       13,
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}))
}

func TestUpdate(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	articleRepository := new(MockNewArticleRepository)
	ctx := context.Background()

	var resp response.Response

	articleRepository.On("Update", ctx, article.Article{
		ID:       1,
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}).Return(nil)
	articleUsecase := article.NewArticleUsecase(nil, location, articleRepository)
	resp = articleUsecase.Update(ctx, article.UpdateArticleRequest{
		ID:       1,
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	})

	assert.Equal(t, resp, response.Success(response.StatusOK, article.UpdateArticleResponses{
		ID:       1,
		Title:    "title1",
		Subtitle: "Indonesia",
		Content:  "Animasi",
	}))
}

func TestDelete(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	articleRepository := new(MockNewArticleRepository)
	ctx := context.Background()

	var resp response.Response

	articleRepository.On("Delete", ctx, int64(1)).Return(nil)
	articleUsecase := article.NewArticleUsecase(nil, location, articleRepository)
	resp = articleUsecase.Delete(ctx, int64(1))
	assert.Equal(t, resp, response.Success(response.StatusOK, nil))
}

func TestSetArticleStatus(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	articleRepository := new(MockNewArticleRepository)
	ctx := context.Background()

	var resp response.Response

	articleRepository.On("SetArticleStatus", ctx, int64(1), "published").Return(nil)
	articleUsecase := article.NewArticleUsecase(nil, location, articleRepository)
	resp = articleUsecase.PublishArticleStatus(ctx, int64(1))
	assert.Equal(t, resp, response.Success(response.StatusOK, nil))
}

func TestSFindByID(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	articleRepository := new(MockNewArticleRepository)
	ctx := context.Background()

	var resp response.Response

	createdAt := time.Now().In(location)

	articleRepository.On("FindByID", ctx, int64(1)).Return(article.Article{
		ID:             1,
		Title:          "title",
		Subtitle:       "Indonesia",
		Content:        "Animasi",
		Status:         "Status",
		CreatedAt:      createdAt,
		PublishedAt:    nil,
		LastModifiedAt: nil,
		Author: account.Account{
			ID: 14,
		},
	}, nil)
	articleUsecase := article.NewArticleUsecase(nil, location, articleRepository)
	resp = articleUsecase.FindByID(ctx, int64(1))

	assert.Equal(t, resp, response.Success(response.StatusOK, article.ArticleResponses{
		ID:             1,
		Title:          "title",
		Subtitle:       "Indonesia",
		Content:        "Animasi",
		Status:         "Status",
		CreatedAt:      createdAt,
		PublishedAt:    nil,
		LastModifiedAt: nil,
		Author: account.Account{
			ID: 14,
		},
	}))
}
