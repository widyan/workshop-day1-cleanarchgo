package article

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sangianpatrick/devoria-article-service/exception"
)

type ArticleRepository interface {
	Save(ctx context.Context, article Article) (ID int64, err error)
	Update(ctx context.Context, updatedArticle Article) (err error)
	Delete(ctx context.Context, ID int64) (err error)
	SetArticleStatus(ctx context.Context, ID int64, status string) (err error)
	FindByID(ctx context.Context, ID int64) (article Article, err error)
	// FindMany(ctx context.Context) (bunchOfArticles []Article, err error)
	// FindManySpecificProfile(ctx context.Context, articleID int64) (bunchOfArticles []Article, err error)
}

type articleRepositoryImpl struct {
	db        *sql.DB
	tableName string
	location  *time.Location
}

func NewArticleRepository(db *sql.DB, tableName string, location *time.Location) ArticleRepository {
	return &articleRepositoryImpl{
		db:        db,
		tableName: tableName,
		location:  location,
	}
}

func (r *articleRepositoryImpl) Save(ctx context.Context, article Article) (ID int64, err error) {
	command := fmt.Sprintf("INSERT INTO %s (authorid, title, subtitle, content, status, createdAt) VALUES (?, ?, ?, ?, ?, ?)", r.tableName)
	stmt, err := r.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		article.Author.ID,
		article.Title,
		article.Subtitle,
		article.Content,
		article.Status,
		time.Now().In(r.location),
	)

	if err != nil {
		log.Println(err)
		return
	}

	ID, _ = result.LastInsertId()

	return
}

func (r *articleRepositoryImpl) Update(ctx context.Context, article Article) (err error) {
	command := fmt.Sprintf("UPDATE %s SET title=?, subtitle=?, content=?, lastModifiedAt=? WHERE id=?", r.tableName)
	stmt, err := r.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		article.Title,
		article.Subtitle,
		article.Content,
		time.Now().In(r.location),
		article.ID,
	)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r *articleRepositoryImpl) Delete(ctx context.Context, ID int64) (err error) {
	command := fmt.Sprintf("DELETE FROM %s WHERE id=?", r.tableName)
	stmt, err := r.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		ID,
	)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r *articleRepositoryImpl) SetArticleStatus(ctx context.Context, ID int64, status string) (err error) {
	command := fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ?`, r.tableName)
	stmt, err := r.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		err = exception.ErrInternalServer
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		status,
		ID,
	)

	if err != nil {
		log.Println(err)
		err = exception.ErrInternalServer
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		err = exception.ErrNotFound
		return
	}

	return
}

func (r *articleRepositoryImpl) FindByID(ctx context.Context, ID int64) (article Article, err error) {
	query := fmt.Sprintf(`SELECT id, authorId, title, subtitle, content, status, createdAt, publishedAt, lastModifiedAt 
	FROM %s WHERE id = ? LIMIT 1`, r.tableName)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, ID)

	err = row.Scan(
		&article.ID,
		&article.Author.ID,
		&article.Title,
		&article.Subtitle,
		&article.Content,
		&article.Status,
		&article.CreatedAt,
		&article.PublishedAt,
		&article.LastModifiedAt,
	)

	if err != nil {
		log.Println(err)
		err = exception.ErrNotFound
		return
	}

	return
}
