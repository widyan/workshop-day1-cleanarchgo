package account

import (
	"encoding/json"
	"fmt"
	"github.com/sangianpatrick/devoria-article-service/entity"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sangianpatrick/devoria-article-service/middleware"
	"github.com/sangianpatrick/devoria-article-service/response"
)

type AccountHTTPHandler struct {
	Validate *validator.Validate
	Usecase  AccountUsecase
}

func NewAccountHTTPHandler(
	router *mux.Router,
	basicAuthMiddleware middleware.RouteMiddleware,
	jwtAuth jwt.JwtMiddleware,
	validate *validator.Validate,
	usecase AccountUsecase,
) {
	handler := &AccountHTTPHandler{
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/v1/account/registration", basicAuthMiddleware.Verify(handler.Register)).Methods(http.MethodPost)
	router.HandleFunc("/v1/account/login", basicAuthMiddleware.Verify(handler.Login)).Methods(http.MethodPost)
	router.HandleFunc("/v1/account", jwtAuth.VerifyToken(handler.Profile)).Methods(http.MethodGet)

}

func (handler *AccountHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var params AccountRegistrationRequest
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

	resp = handler.Usecase.Register(ctx, params)
	resp.JSON(w)
}

func (handler *AccountHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var params AccountAuthenticationRequest
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

	resp = handler.Usecase.Login(ctx, params)
	resp.JSON(w)
}

func (handler *AccountHTTPHandler) Profile(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var ctx = r.Context()

	var err error
	var claims entity.AccountStandardJWTClaims
	bind, ok := context.Get(r,"bind").([]byte)
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
	resp = handler.Usecase.GetProfile(ctx, claims)
	resp.JSON(w)
}
