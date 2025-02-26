package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	log "go.uber.org/zap"

	"timble/internal/utils"
	mockshandler "timble/mocks/module/users/internal_/usecase"
	"timble/module/users/entity"
	"timble/module/users/internal/handler"
	"timble/module/users/internal/repository"
	"timble/module/users/internal/usecase"
)

var (
	normalTokenResponseData = entity.UserToken{
		Token: "aaaaaaaaaa.aaaaaaaaaaa.aaaaa-aaaaa-aaaaa_aaaaaa",
	}

	normalTokenResponseString = `{
      "meta":{
         "http_status":%d
      },
      "data": {
         "token": "aaaaaaaaaa.aaaaaaaaaaa.aaaaa-aaaaa-aaaaa_aaaaaa"
      }
   }`

	stdErrorResponseBase = `{
   	"meta":{
      	"http_status":%d
   	},
   	"error_detail": {
   		"message": "%s",
   		"code": "%s",
   		"field": "%s"
   	}
	}`

	stdErrorResponseWithoutField = `{
   	"meta":{
      	"http_status":%d
   	},
   	"error_detail": {
   		"message": "%s",
   		"code": "%s"
   	}
	}`

	messageResponseBase = `{
   	"meta":{
      	"http_status":%d
   	},
   	"message":"%s"
	}`
)

type shouldMock struct {
	handlerFunc bool
}

type expected struct {
	expectedResponse   string
	expectedHTTPStatus int
}

func Test_NewUsersResource(t *testing.T) {
	t.Run("new users resource", func(t *testing.T) {
		auc := usecase.NewAuthUsecase(&utils.AuthConfig{}, &repository.PostgresRepository{}, &log.Logger{})
		puc := usecase.NewPremiumUsecase(&repository.RedisRepository{}, &repository.PostgresRepository{}, &repository.CacheRepository{}, &log.Logger{})
		uuc := usecase.NewUserUsecase(&utils.AuthConfig{}, &repository.RedisRepository{}, &repository.PostgresRepository{}, &repository.CacheRepository{}, &log.Logger{})

		res := handler.NewUsersResource(auc, puc, uuc, &log.Logger{})

		assert.IsType(t, &handler.UsersResource{}, res)
	})
}

func TestUsersResource_Login(t *testing.T) {
	normalRequestData := `{
      "username":  "testuser",
      "password": "testpassword"
    }`

	normalRequestDataParsed := entity.UserLoginParams{
		Username: "testuser",
		Password: "testpassword",
	}

	badRequestData := `{
      "username":  "testuser",
      "password": ""
    }`

	type args struct {
		requestData       string
		requestDataParsed entity.UserLoginParams
	}

	type mocked struct {
		handlerResult entity.UserToken
		handlerError  error
	}

	cases := []struct {
		name       string
		args       args
		mocked     mocked
		shouldMock shouldMock
		expected   expected
	}{
		{
			name: "normal case - successfully retrieve token",
			args: args{
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerResult: normalTokenResponseData,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusOK,
				expectedResponse:   fmt.Sprintf(normalTokenResponseString, http.StatusOK),
			},
		},
		{
			name: "error case - missing parameters",
			args: args{
				requestData: badRequestData,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "Password must be more than 10 characters", "PARAMETER_PARSING_FAILS", "password"),
			},
		},
		{
			name: "error case - handler returned standard error",
			args: args{
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: utils.ErrorInvalidLogin,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusUnauthorized,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusUnauthorized, "Invalid username or password", "Unauthorized"),
			},
		},
		{
			name: "error case - handler returned unexpected error",
			args: args{
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: errors.New("unexpected"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusInternalServerError,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusInternalServerError, "internal server error, please check the server logs", "INTERNAL SERVER ERROR"),
			},
		},
	}

	for _, tc := range cases {
		uc := mockshandler.NewAuthUsecase(t)
		t.Run(tc.name, func(t *testing.T) {
			logger := initLogger(t)
			urlPath := "/api/public/auth/login"

			req := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer([]byte(tc.args.requestData)))
			recorder := httptest.NewRecorder()
			ctx := initRoutingContext(req.Context())
			req = req.WithContext(ctx)

			if tc.shouldMock.handlerFunc {
				uc.
					On("Login", ctx, tc.args.requestDataParsed).
					Return(tc.mocked.handlerResult, tc.mocked.handlerError)
			}

			st := handler.NewUsersResource(uc, mockshandler.NewPremiumUsecase(t), mockshandler.NewUserUsecase(t), logger)

			hndlr := http.HandlerFunc(st.Login)
			hndlr.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
			assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
		})
	}
}

func TestUsersResource_Create(t *testing.T) {
	normalRequestData := `{
      "username":  "testuser",
      "email": "test@email.com",
      "password": "testpassword"
    }`

	normalRequestDataParsed := entity.UserRegistrationParams{
		Username: "testuser",
		Email:    "test@email.com",
		Password: "testpassword",
	}

	badRequestData := `{
      "username":  "testuser",
      "password": "testpassword"
    }`

	type args struct {
		requestData       string
		requestDataParsed entity.UserRegistrationParams
	}

	type mocked struct {
		handlerResult entity.UserToken
		handlerError  error
	}

	cases := []struct {
		name       string
		args       args
		mocked     mocked
		shouldMock shouldMock
		expected   expected
	}{
		{
			name: "normal case - successfully create user and retrieve token",
			args: args{
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerResult: normalTokenResponseData,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusCreated,
				expectedResponse:   fmt.Sprintf(normalTokenResponseString, http.StatusCreated),
			},
		},
		{
			name: "error case - missing parameters",
			args: args{
				requestData: badRequestData,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "Email can not be blank", "PARAMETER_PARSING_FAILS", "email"),
			},
		},
		{
			name: "error case - handler returned standard error",
			args: args{
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: utils.NewStandardError("unexpected", "DB ERROR", "server"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "unexpected", "DB ERROR", "server"),
			},
		},
		{
			name: "error case - handler returned unexpected error",
			args: args{
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: errors.New("unexpected"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusInternalServerError,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusInternalServerError, "internal server error, please check the server logs", "INTERNAL SERVER ERROR"),
			},
		},
	}

	for _, tc := range cases {
		uc := mockshandler.NewUserUsecase(t)
		t.Run(tc.name, func(t *testing.T) {
			logger := initLogger(t)
			urlPath := "/api/public/user/register"

			req := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer([]byte(tc.args.requestData)))
			recorder := httptest.NewRecorder()
			ctx := initRoutingContext(req.Context())
			req = req.WithContext(ctx)

			if tc.shouldMock.handlerFunc {
				uc.
					On("Create", ctx, tc.args.requestDataParsed).
					Return(tc.mocked.handlerResult, tc.mocked.handlerError)
			}

			st := handler.NewUsersResource(mockshandler.NewAuthUsecase(t), mockshandler.NewPremiumUsecase(t), uc, logger)

			hndlr := http.HandlerFunc(st.Create)
			hndlr.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
			assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
		})
	}
}

func TestUsersResource_Show(t *testing.T) {
	timestamp, _ := time.Parse("1/2/2006", "2/2/2025")
	normalUser := &entity.UserPublic{
		ID:        uint(1),
		Email:     "test@email.com",
		Username:  "testuser",
		Premium:   true,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	normalUserResponseString := `{
	   "meta":{
	      "http_status":200
	   },
	   "data":{
	      "id":1,
	      "email":"test@email.com",
	      "username":"testuser",
	      "premium":true,
	      "created_at":"2025-02-02T00:00:00Z",
	      "updated_at":"2025-02-02T00:00:00Z"
	   }
	}`

	type args struct {
		args uint
	}

	type mocked struct {
		handlerResult *entity.UserPublic
		handlerError  error
	}

	cases := []struct {
		name       string
		args       args
		mocked     mocked
		shouldMock shouldMock
		expected   expected
	}{
		{
			name: "normal case - successfully get user",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerResult: normalUser,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusOK,
				expectedResponse:   normalUserResponseString,
			},
		},
		{
			name: "error case - user not found",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusBadRequest, "User not found:1", "NOT FOUND"),
			},
		},
		{
			name: "error case - handler returned standard error",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: utils.NewStandardError("unexpected", "DB ERROR", "server"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "unexpected", "DB ERROR", "server"),
			},
		},
		{
			name: "error case - handler returned unexpected error",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: errors.New("unexpected"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusInternalServerError,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusInternalServerError, "internal server error, please check the server logs", "INTERNAL SERVER ERROR"),
			},
		},
	}

	for _, tc := range cases {
		uc := mockshandler.NewUserUsecase(t)
		t.Run(tc.name, func(t *testing.T) {
			logger := initLogger(t)
			urlPath := "/api/protected/users"

			req := httptest.NewRequest(http.MethodGet, urlPath, bytes.NewBuffer(nil))
			recorder := httptest.NewRecorder()
			ctx := initRoutingContext(context.WithValue(req.Context(), utils.CtxUserIDKey, float64(tc.args.args)))
			req = req.WithContext(ctx)

			if tc.shouldMock.handlerFunc {
				uc.
					On("Show", ctx, tc.args.args).
					Return(tc.mocked.handlerResult, tc.mocked.handlerError)
			}

			st := handler.NewUsersResource(mockshandler.NewAuthUsecase(t), mockshandler.NewPremiumUsecase(t), uc, logger)

			hndlr := http.HandlerFunc(st.Show)
			hndlr.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
			assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
		})
	}
}

func TestUsersResource_React(t *testing.T) {
	normalRequestData := `{
		      "target_id":  2,
		      "type": 1
		    }`

	normalRequestDataParsed := entity.ReactionParams{
		UserID:   1,
		TargetID: 2,
		Type:     1,
	}

	badRequestData := `{
		      "type": 1
		    }`

	type args struct {
		args              uint
		requestData       string
		requestDataParsed entity.ReactionParams
	}

	type mocked struct {
		handlerResult string
		handlerError  error
	}

	cases := []struct {
		name       string
		args       args
		mocked     mocked
		shouldMock shouldMock
		expected   expected
	}{
		{
			name: "normal case - successfully react",
			args: args{
				args:              1,
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerResult: "Reaction saved",
			},
			expected: expected{
				expectedHTTPStatus: http.StatusOK,
				expectedResponse:   fmt.Sprintf(messageResponseBase, http.StatusOK, "Reaction saved"),
			},
		},
		{
			name: "error case - missing parameters",
			args: args{
				args:        1,
				requestData: badRequestData,
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "Invalid target user", "PARAMETER_PARSING_FAILS", "target_id"),
			},
		},
		{
			name: "error case - handler returned standard error",
			args: args{
				args:              1,
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: utils.NewStandardError("unexpected", "DB ERROR", "server"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "unexpected", "DB ERROR", "server"),
			},
		},
		{
			name: "error case - handler returned unexpected error",
			args: args{
				args:              1,
				requestData:       normalRequestData,
				requestDataParsed: normalRequestDataParsed,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: errors.New("unexpected"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusInternalServerError,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusInternalServerError, "internal server error, please check the server logs", "INTERNAL SERVER ERROR"),
			},
		},
	}

	for _, tc := range cases {
		uc := mockshandler.NewUserUsecase(t)
		t.Run(tc.name, func(t *testing.T) {
			logger := initLogger(t)
			urlPath := "/api/protected/users/react"

			req := httptest.NewRequest(http.MethodPatch, urlPath, bytes.NewBuffer([]byte(tc.args.requestData)))
			recorder := httptest.NewRecorder()
			ctx := initRoutingContext(context.WithValue(req.Context(), utils.CtxUserIDKey, float64(tc.args.args)))
			req = req.WithContext(ctx)

			if tc.shouldMock.handlerFunc {
				uc.
					On("React", ctx, tc.args.requestDataParsed).
					Return(tc.mocked.handlerError)
			}

			st := handler.NewUsersResource(mockshandler.NewAuthUsecase(t), mockshandler.NewPremiumUsecase(t), uc, logger)

			hndlr := http.HandlerFunc(st.React)
			hndlr.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
			assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
		})
	}
}

func TestUsersResource_GrantPremium(t *testing.T) {
	type args struct {
		args uint
	}

	type mocked struct {
		handlerResult string
		handlerError  error
	}

	cases := []struct {
		name       string
		args       args
		mocked     mocked
		shouldMock shouldMock
		expected   expected
	}{
		{
			name: "normal case - successfully grant premium",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerResult: "Reaction saved",
			},
			expected: expected{
				expectedHTTPStatus: http.StatusOK,
				expectedResponse:   fmt.Sprintf(messageResponseBase, http.StatusOK, "Premium granted"),
			},
		},
		{
			name: "error case - handler returned standard error",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: utils.NewStandardError("unexpected", "DB ERROR", "server"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "unexpected", "DB ERROR", "server"),
			},
		},
		{
			name: "error case - handler returned unexpected error",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: errors.New("unexpected"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusInternalServerError,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusInternalServerError, "internal server error, please check the server logs", "INTERNAL SERVER ERROR"),
			},
		},
	}

	for _, tc := range cases {
		uc := mockshandler.NewPremiumUsecase(t)
		t.Run(tc.name, func(t *testing.T) {
			logger := initLogger(t)
			urlPath := "/api/protected/users/premium/grant"

			req := httptest.NewRequest(http.MethodPatch, urlPath, bytes.NewBuffer(nil))
			recorder := httptest.NewRecorder()
			ctx := initRoutingContext(context.WithValue(req.Context(), utils.CtxUserIDKey, float64(tc.args.args)))
			req = req.WithContext(ctx)

			if tc.shouldMock.handlerFunc {
				uc.
					On("Grant", ctx, tc.args.args).
					Return(tc.mocked.handlerError)
			}

			st := handler.NewUsersResource(mockshandler.NewAuthUsecase(t), uc, mockshandler.NewUserUsecase(t), logger)

			hndlr := http.HandlerFunc(st.GrantPremium)
			hndlr.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
			assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
		})
	}
}

func TestUsersResource_UnsubscribePremium(t *testing.T) {
	type args struct {
		args uint
	}

	type mocked struct {
		handlerResult string
		handlerError  error
	}

	cases := []struct {
		name       string
		args       args
		mocked     mocked
		shouldMock shouldMock
		expected   expected
	}{
		{
			name: "normal case - successfully unsubscribed from premium",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerResult: "Reaction saved",
			},
			expected: expected{
				expectedHTTPStatus: http.StatusOK,
				expectedResponse:   fmt.Sprintf(messageResponseBase, http.StatusOK, "Unsubscribed from premium"),
			},
		},
		{
			name: "error case - handler returned standard error",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: utils.NewStandardError("unexpected", "DB ERROR", "server"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusBadRequest,
				expectedResponse:   fmt.Sprintf(stdErrorResponseBase, http.StatusBadRequest, "unexpected", "DB ERROR", "server"),
			},
		},
		{
			name: "error case - handler returned unexpected error",
			args: args{
				args: 1,
			},
			shouldMock: shouldMock{
				handlerFunc: true,
			},
			mocked: mocked{
				handlerError: errors.New("unexpected"),
			},
			expected: expected{
				expectedHTTPStatus: http.StatusInternalServerError,
				expectedResponse:   fmt.Sprintf(stdErrorResponseWithoutField, http.StatusInternalServerError, "internal server error, please check the server logs", "INTERNAL SERVER ERROR"),
			},
		},
	}

	for _, tc := range cases {
		uc := mockshandler.NewPremiumUsecase(t)
		t.Run(tc.name, func(t *testing.T) {
			logger := initLogger(t)
			urlPath := "/api/protected/users/premium/unsubscribe"

			req := httptest.NewRequest(http.MethodPatch, urlPath, bytes.NewBuffer(nil))
			recorder := httptest.NewRecorder()
			ctx := initRoutingContext(context.WithValue(req.Context(), utils.CtxUserIDKey, float64(tc.args.args)))
			req = req.WithContext(ctx)

			if tc.shouldMock.handlerFunc {
				uc.
					On("Unsubscribe", ctx, tc.args.args).
					Return(tc.mocked.handlerError)
			}

			st := handler.NewUsersResource(mockshandler.NewAuthUsecase(t), uc, mockshandler.NewUserUsecase(t), logger)

			hndlr := http.HandlerFunc(st.UnsubscribePremium)
			hndlr.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.expectedHTTPStatus, recorder.Result().StatusCode)
			assert.JSONEq(t, tc.expected.expectedResponse, recorder.Body.String())
		})
	}
}

func initRoutingContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, chi.RouteCtxKey, chi.NewRouteContext())
}

func initLogger(t *testing.T) *log.Logger {
	config := log.NewProductionConfig()
	config.Level = log.NewAtomicLevelAt(log.FatalLevel)

	logger, err := config.Build()
	assert.Nil(t, err)
	return logger
}
