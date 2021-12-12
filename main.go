package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sangianpatrick/devoria-article-service/config"
	"github.com/sangianpatrick/devoria-article-service/crypto"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/domain/article"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/middleware"
	"github.com/sangianpatrick/devoria-article-service/session"
)

func main() {
	location, _ := time.LoadLocation("Asia/Jakarta")
	cfg := config.New()

	db, err := sql.Open("mysql", cfg.Mariadb.DSN)
	db.SetMaxOpenConns(cfg.Mariadb.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.Mariadb.MaxIdleConnections)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rc := redis.NewClient(cfg.Redis.Options)
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}

	vld := validator.New()
	encryption := crypto.NewAES256CBC(cfg.AES.SecretKey)
	jsonWebToken := jwt.NewJSONWebToken(jwt.GetRSAPrivateKey("./secret/id_rsa"), jwt.GetRSAPublicKey("./secret/id_rsa.pub"))
	sess := session.NewRedisSessionStoreAdapter(rc, time.Hour*24*1)
	basicAuthMiddleware := middleware.NewBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)
	jwtAuthMiddleware := jwt.NewJwtToken(jsonWebToken)

	router := mux.NewRouter()

	accountRepository := account.NewAccountRepository(db, "account")
	accountUsecase := account.NewAccountUsecase(cfg.GlobalIV, sess, jsonWebToken, encryption, location, accountRepository)
	account.NewAccountHTTPHandler(router, basicAuthMiddleware, jwtAuthMiddleware, vld, accountUsecase)

	articleRepository := article.NewArticleRepository(db, "article", location)
	articleUsecase := article.NewArticleUsecase(sess, location, articleRepository)
	article.NewAccountHTTPHandler(router, basicAuthMiddleware, jwtAuthMiddleware, vld, articleUsecase)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	fmt.Println("shutting down application ...")

	server.Shutdown(context.Background())
	db.Close()
	rc.Close()
}
