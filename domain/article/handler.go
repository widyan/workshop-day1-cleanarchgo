package article

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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
	validate *validator.Validate,
	usecase ArticleUsecase,
) {
	handler := &AccountHTTPHandler{
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/v1/article/create", handler.Save).Methods(http.MethodPost)
	router.HandleFunc("/v1/article/update", handler.Update).Methods(http.MethodPut)
	router.HandleFunc("/v1/article/delete/{id}", handler.Delete).Methods(http.MethodDelete)

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

	resp = handler.Usecase.Save(ctx, 13, params)
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
