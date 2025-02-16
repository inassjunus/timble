package handler

import (
	"net/http"

	"github.com/pkg/errors"

	log "go.uber.org/zap"

	"timble/internal/utils"
	"timble/module/users/entity"
	"timble/module/users/internal/usecase"
)

type UsersResource struct {
	AuthUsecase    usecase.AuthUsecase
	PremiumUsecase usecase.PremiumUsecase
	UserUsecase    usecase.UserUsecase
	logger         *log.Logger
}

func NewUsersResource(authUsecase usecase.AuthUsecase, premiumUsecase usecase.PremiumUsecase, userUsecase usecase.UserUsecase, logger *log.Logger) *UsersResource {
	return &UsersResource{
		AuthUsecase:    authUsecase,
		PremiumUsecase: premiumUsecase,
		UserUsecase:    userUsecase,
		logger:         logger,
	}
}

func (resource *UsersResource) Login(w http.ResponseWriter, r *http.Request) {
	m := utils.NewRestMetric(r)

	defer func() {
		m.TrackRestService()
	}()

	params, err := entity.NewUserLoginPayload(r.Body)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	result, err := resource.AuthUsecase.Login(r.Context(), params)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	m.HTTPStatus = http.StatusOK
	meta := utils.Meta{
		HTTPStatus: http.StatusOK,
	}
	body := utils.NewDataResponse(result, meta)
	body.WriteAPIResponse(w, r, http.StatusOK)
}

func (resource *UsersResource) Create(w http.ResponseWriter, r *http.Request) {
	m := utils.NewRestMetric(r)

	defer func() {
		m.TrackRestService()
	}()

	params, err := entity.NewUserRegistrationPayload(r.Body)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	result, err := resource.UserUsecase.Create(r.Context(), params)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	m.HTTPStatus = http.StatusCreated
	meta := utils.Meta{
		HTTPStatus: http.StatusCreated,
	}
	body := utils.NewDataResponse(result, meta)
	body.WriteAPIResponse(w, r, http.StatusCreated)
}

func (resource *UsersResource) Show(w http.ResponseWriter, r *http.Request) {
	userID := resource.getUserIDFromContext(r)
	m := utils.NewRestMetric(r)

	defer func() {
		m.TrackRestService()
	}()

	userData, err := resource.UserUsecase.Show(r.Context(), userID)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	if userData == nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, utils.UserNotFoundError(userID)))
		return
	}

	meta := utils.Meta{
		HTTPStatus: http.StatusOK,
	}
	body := utils.NewDataResponse(userData, meta)
	body.WriteAPIResponse(w, r, http.StatusOK)
}

func (resource *UsersResource) React(w http.ResponseWriter, r *http.Request) {
	userID := resource.getUserIDFromContext(r)
	m := utils.NewRestMetric(r)

	defer func() {
		m.TrackRestService()
	}()

	params, err := entity.NewReactionPayload(r.Body, userID)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	err = resource.UserUsecase.React(r.Context(), params)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	m.HTTPStatus = http.StatusOK
	meta := utils.Meta{
		HTTPStatus: http.StatusOK,
	}
	body := utils.NewMessageResponse("Reaction saved", meta)
	body.WriteAPIResponse(w, r, http.StatusOK)
}

func (resource *UsersResource) GrantPremium(w http.ResponseWriter, r *http.Request) {
	userID := resource.getUserIDFromContext(r)
	m := utils.NewRestMetric(r)

	defer func() {
		m.TrackRestService()
	}()

	err := resource.PremiumUsecase.Grant(r.Context(), userID)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	meta := utils.Meta{
		HTTPStatus: http.StatusOK,
	}
	body := utils.NewMessageResponse("Premium granted", meta)
	body.WriteAPIResponse(w, r, http.StatusOK)
}

func (resource *UsersResource) UnsubscribePremium(w http.ResponseWriter, r *http.Request) {
	userID := resource.getUserIDFromContext(r)
	m := utils.NewRestMetric(r)

	defer func() {
		m.TrackRestService()
	}()

	err := resource.PremiumUsecase.Unsubscribe(r.Context(), userID)
	if err != nil {
		m = m.SetFail(resource.returnErrorResponse(w, r, err))
		return
	}

	meta := utils.Meta{
		HTTPStatus: http.StatusOK,
	}
	body := utils.NewMessageResponse("Unsubscribed from premium", meta)
	body.WriteAPIResponse(w, r, http.StatusOK)
}

func (resource *UsersResource) getUserIDFromContext(r *http.Request) uint {
	return uint(r.Context().Value("user_id").(float64))
}

func (resource *UsersResource) returnErrorResponse(w http.ResponseWriter, r *http.Request, err error) int {
	errOrig, ok := err.(*utils.StandardError)
	if !ok {
		httpStatus := http.StatusInternalServerError
		resource.logger.Error(errors.WithStack(err).Error(), utils.BuildRequestLogFields(r, httpStatus)...)
		writeErrorResponse(w, r, utils.ErrorInternalServerResponse, httpStatus)
		return httpStatus
	}

	httpStatus := http.StatusBadRequest
	if errOrig.HttpStatus != 0 {
		httpStatus = errOrig.HttpStatus
	}
	writeErrorResponse(w, r, errOrig, httpStatus)
	resource.logger.Error(err.Error(), utils.BuildRequestLogFields(r, httpStatus)...)
	return httpStatus
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, err *utils.StandardError, statusCode int) {
	body := utils.NewErrorResponse(err, statusCode)
	body.WriteAPIResponse(w, r, statusCode)
}
