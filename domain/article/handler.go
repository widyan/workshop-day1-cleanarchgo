package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sangianpatrick/devoria-article-service/entity"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/middleware"
	"github.com/sangianpatrick/devoria-article-service/response"
)

type AccountHTTPHandler struct {
	Validate *validator.Validate
	Usecase  ArticleUsecase
}

func NewAccountHTTPHandler(
	router *mux.Router,
	basicAuthMiddleware middleware.RouteMiddleware,
	jwtAuth jwt.JwtMiddleware,
	validate *validator.Validate,
	usecase ArticleUsecase,
) {
	handler := &AccountHTTPHandler{
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/v1/article/create", jwtAuth.VerifyToken(handler.Save)).Methods(http.MethodPost)
	router.HandleFunc("/v1/article/update", jwtAuth.VerifyToken(handler.Update)).Methods(http.MethodPut)
	router.HandleFunc("/v1/article/delete/{id}", jwtAuth.VerifyToken(handler.Delete)).Methods(http.MethodDelete)
	router.HandleFunc("/v1/article/publish/{id}", jwtAuth.VerifyToken(handler.PublishArticleStatus)).Methods(http.MethodPut)
	router.HandleFunc("/v1/article/findbyid/{id}", jwtAuth.VerifyToken(handler.FindByID)).Methods(http.MethodGet)
}

func (handler *AccountHTTPHandler) Save(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var params CreateArticleRequest
	var ctx = r.Context()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		resp = response.Error(response.StatusUnprocessabelEntity, nil, err)
		resp.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, params)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}

	var claims entity.AccountStandardJWTClaims
	bind, ok := context.Get(r, "bind").([]byte)
	if !ok {
		err = fmt.Errorf("Error Bind Value")
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}

	err = json.Unmarshal(bind, &claims)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}

	id, err := strconv.Atoi(claims.StandardClaims.Subject)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}

	resp = handler.Usecase.Save(ctx, int64(id), params)
	resp.JSON(w)
}

func (handler *AccountHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var params UpdateArticleRequest
	var ctx = r.Context()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		resp = response.Error(response.StatusUnprocessabelEntity, nil, err)
		resp.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, params)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}

	resp = handler.Usecase.Update(ctx, params)
	resp.JSON(w)
}

func (handler *AccountHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {

	var resp response.Response

	var ctx = r.Context()
	ids := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}
	resp = handler.Usecase.Delete(ctx, id)
	resp.JSON(w)
}

func (handler *AccountHTTPHandler) PublishArticleStatus(w http.ResponseWriter, r *http.Request) {

	var resp response.Response

	var ctx = r.Context()
	ids := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}
	resp = handler.Usecase.PublishArticleStatus(ctx, id)
	resp.JSON(w)
}

func (handler *AccountHTTPHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	var resp response.Response

	var ctx = r.Context()
	ids := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}
	resp = handler.Usecase.FindByID(ctx, id)
	resp.JSON(w)
}
