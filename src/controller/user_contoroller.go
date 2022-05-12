package controller

import (
	"context"
	"net/http"
	"url-shortener/config"
	"url-shortener/src/models/url_shortener"
	"url-shortener/src/models/users"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UsrUsecase users.UserUsecase
	UrlUsecase url_shortener.UrlUsecase
}

func NewUserHandler(e *echo.Echo, userJwt *echo.Group, uu users.UserUsecase, usecase url_shortener.UrlUsecase, configuration *config.Config) {
	handler := &UserHandler{
		UsrUsecase: uu,
		UrlUsecase: usecase,
	}
	e.GET("/:key", handler.Redirect)
	e.POST("/user", handler.InsertOne)
	userJwt.GET("/user", handler.FindOne)
	userJwt.PUT("/user", handler.UpdateOne)
}

// check is valid request
func isRequestValid(m *users.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Redirect this end point is for redirect to main url
//@param is key for getting original url
func (user *UserHandler) Redirect(c echo.Context) error {
	//get key form main url
	url := c.Param("key")
	ctx := c.Request().Context()
	// get original url by key
	OriginalURL, err := user.UrlUsecase.FindOneByKey(ctx, url)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	//redirect ro original url
	return c.Redirect(http.StatusMovedPermanently, OriginalURL)
}

// InsertOne create user
func (user *UserHandler) InsertOne(c echo.Context) error {
	var (
		usr users.User
		err error
	)
	// get user paramter
	err = c.Bind(&usr)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&usr); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	//insert user to mongo database
	result, err := user.UsrUsecase.InsertOne(ctx, &usr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// FindOne TODO list FindOne
func (user *UserHandler) FindOne(c echo.Context) error {

	id := c.QueryParam("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.FindOne(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateOne TODO list  UpdateOne
func (user *UserHandler) UpdateOne(c echo.Context) error {

	id := c.QueryParam("id")

	var (
		usr users.User
		err error
	)

	err = c.Bind(&usr)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.UpdateOne(ctx, &usr, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
