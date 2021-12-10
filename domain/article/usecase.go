package article

import (
	"context"
	"time"

	"github.com/sangianpatrick/devoria-article-service/crypto"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/exception"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/response"
	"github.com/sangianpatrick/devoria-article-service/session"
)

type ArticleUsecase interface {
	Save(ctx context.Context, request CreateArticleRequest) (resp response.Response)
	Update(ctx context.Context, request UpdateArticleRequest) (resp response.Response)
	Delete(ctx context.Context, ID int64) (resp response.Response)
}

type articleUsecaseImpl struct {
	globalIV     string
	session      session.Session
	jsonWebToken jwt.JSONWebToken
	crypto       crypto.Crypto
	location     *time.Location
	repository   ArticleRepository
}

func NewArticleUsecase(
	globalIV string,
	session session.Session,
	jsonWebToken jwt.JSONWebToken,
	crypto crypto.Crypto,
	location *time.Location,
	repository ArticleRepository,
) ArticleUsecase {
	return &articleUsecaseImpl{
		globalIV:     globalIV,
		session:      session,
		jsonWebToken: jsonWebToken,
		crypto:       crypto,
		location:     location,
		repository:   repository,
	}
}

func (u *articleUsecaseImpl) Save(ctx context.Context, article CreateArticleRequest) (resp response.Response) {
	createdAt := time.Now().In(u.location)
	_, err := u.repository.Save(ctx, Article{
		Title:     article.Title,
		Subtitle:  article.Subtitle,
		Content:   article.Content,
		CreatedAt: createdAt,
		Author: account.Account{
			ID: 13,
		},
	})
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	return response.Success(response.StatusCreated, article)
}

func (u *articleUsecaseImpl) Update(ctx context.Context, article UpdateArticleRequest) (resp response.Response) {
	lastModifiedAt := time.Now().In(u.location)
	err := u.repository.Update(ctx, Article{
		ID:             article.ID,
		Title:          article.Title,
		Subtitle:       article.Subtitle,
		Content:        article.Content,
		LastModifiedAt: &lastModifiedAt,
	})
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, article)
}

func (u *articleUsecaseImpl) Delete(ctx context.Context, ID int64) (resp response.Response) {
	err := u.repository.Delete(ctx, ID)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, nil)
}
