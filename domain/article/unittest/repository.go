package unittest

import (
	"context"

	"github.com/sangianpatrick/devoria-article-service/domain/article"
	"github.com/stretchr/testify/mock"
)

type MockNewArticleRepository struct {
	mock.Mock
}

func (d *MockNewArticleRepository) Save(ctx context.Context, articles article.Article) (ID int64, err error) {
	args := d.Called(ctx, articles)
	return int64(args.Int(0)), args.Error(1)
}

func (d *MockNewArticleRepository) Update(ctx context.Context, updatedArticle article.Article) (err error) {
	args := d.Called(ctx, updatedArticle)
	return args.Error(0)
}

func (d *MockNewArticleRepository) Delete(ctx context.Context, ID int64) (err error) {
	args := d.Called(ctx, ID)
	return args.Error(0)
}

func (d *MockNewArticleRepository) SetArticleStatus(ctx context.Context, ID int64, status string) (err error) {
	args := d.Called(ctx, ID, status)
	return args.Error(0)
}

func (d *MockNewArticleRepository) FindByID(ctx context.Context, ID int64) (articles article.Article, err error) {
	args := d.Called(ctx, ID)
	return args.Get(0).(article.Article), args.Error(1)
}
