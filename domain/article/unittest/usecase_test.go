package unittest

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/sangianpatrick/devoria-article-service/config"
	"github.com/sangianpatrick/devoria-article-service/crypto"
	"github.com/sangianpatrick/devoria-article-service/domain/article"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/response"
	"github.com/stretchr/testify/assert"
)

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
