package article

import (
	"context"
	"time"

	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/exception"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/response"
	"github.com/sangianpatrick/devoria-article-service/session"
)

type ArticleUsecase interface {
	Save(ctx context.Context, AuthorID int64, request CreateArticleRequest) (resp response.Response)
	Update(ctx context.Context, request UpdateArticleRequest) (resp response.Response)
	Delete(ctx context.Context, ID int64) (resp response.Response)
	PublishArticleStatus(ctx context.Context, articleID int64) (resp response.Response)
	FindByID(ctx context.Context, articleID int64) (resp response.Response)
}

type articleUsecaseImpl struct {
	globalIV     string
	session      session.Session
	jsonWebToken jwt.JSONWebToken
	location     *time.Location
	repository   ArticleRepository
}

func NewArticleUsecase(
	session session.Session,
	location *time.Location,
	repository ArticleRepository,
) ArticleUsecase {
	return &articleUsecaseImpl{
		session:    session,
		location:   location,
		repository: repository,
	}
}

func (u *articleUsecaseImpl) Save(ctx context.Context, AuthorID int64, article CreateArticleRequest) (resp response.Response) {
	id, err := u.repository.Save(ctx, Article{
		Title:    article.Title,
		Subtitle: article.Subtitle,
		Content:  article.Content,
		Author: account.Account{
			ID: AuthorID,
		},
	})
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	return response.Success(response.StatusCreated, CreateArticleResponses{
		ID:       id,
		Title:    article.Title,
		Subtitle: article.Subtitle,
		Content:  article.Content,
	})
}

func (u *articleUsecaseImpl) Update(ctx context.Context, article UpdateArticleRequest) (resp response.Response) {
	err := u.repository.Update(ctx, Article{
		ID:       article.ID,
		Title:    article.Title,
		Subtitle: article.Subtitle,
		Content:  article.Content,
	})
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, UpdateArticleResponses{
		ID:       article.ID,
		Title:    article.Title,
		Subtitle: article.Subtitle,
		Content:  article.Content,
	})
}

func (u *articleUsecaseImpl) Delete(ctx context.Context, ID int64) (resp response.Response) {
	err := u.repository.Delete(ctx, ID)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, nil)
}

func (u *articleUsecaseImpl) PublishArticleStatus(ctx context.Context, articleID int64) (resp response.Response) {
	err := u.repository.SetArticleStatus(ctx, articleID, "published")
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}
	return response.Success(response.StatusOK, nil)
}

func (u *articleUsecaseImpl) FindByID(ctx context.Context, articleID int64) (resp response.Response) {
	article, err := u.repository.FindByID(ctx, articleID)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}
	return response.Success(response.StatusOK, ArticleResponses{
		ID:             article.ID,
		Title:          article.Title,
		Subtitle:       article.Subtitle,
		Content:        article.Content,
		Status:         article.Status,
		CreatedAt:      article.CreatedAt,
		PublishedAt:    article.PublishedAt,
		LastModifiedAt: article.LastModifiedAt,
		Author: account.Account{
			ID: article.Author.ID,
		},
	})
}
