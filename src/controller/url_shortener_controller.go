package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"url-shortener/config"
	"url-shortener/src/models/url_shortener"
	"url-shortener/src/models/users"
)

type ResponseUrlError struct {
	Message string `json:"message"`
}

type UrlShotenerHandler struct {
	UrlShotenerUsecase url_shortener.UrlUsecase
	Database           *config.MemoryClient
	configuration      *config.Config
}

func NewUrlHandler(userJwt *echo.Group, uu url_shortener.UrlUsecase, database *config.MemoryClient, configuration *config.Config) {
	handler := &UrlShotenerHandler{
		UrlShotenerUsecase: uu,
		Database:           database,
		configuration:      configuration,
	}
	userJwt.POST("/UrlShotener", handler.InsertOne)
	//userJwt.GET("/UrlShotener", handler.FindOne)
	//
	//userJwt.PUT("/UrlShotener", handler.UpdateOne)
}

// check is valid paramter
func isRequestUrlValid(m *url_shortener.UrlShortener) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// InsertOne create short url for Insert one
func (UrlShotener *UrlShotenerHandler) InsertOne(c echo.Context) error {
	var (
		ct  url_shortener.UrlShortener
		err error
	)
	//get user from context
	user := c.Get("user")
	token := user.(*users.User)

	// bind params to url_shortener.UrlShortener
	err = c.Bind(&ct)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	//check is valid request
	if ok, err := isRequestUrlValid(&ct); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ct.UserID = token.ID
	// insert to redis and main
	result, err := UrlShotener.UrlShotenerUsecase.InsertOne(ctx, &ct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result.ShortUrl = fmt.Sprintf("%s%s", UrlShotener.configuration.BaseUrl, result.ShortUrl)

	return c.JSON(http.StatusOK, result)
}

// FindOne TODO list FindOne
func (UrlShotener *UrlShotenerHandler) FindOne(c echo.Context) error {

	id := c.QueryParam("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := UrlShotener.UrlShotenerUsecase.FindOne(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result.ShortUrl = fmt.Sprintf("%s%s", UrlShotener.configuration.BaseUrl, result.ShortUrl)
	return c.JSON(http.StatusOK, result)
}

// UpdateOne TODO list UpdateOne
func (UrlShotener *UrlShotenerHandler) UpdateOne(c echo.Context) error {

	id := c.QueryParam("id")

	var (
		ct  url_shortener.UrlShortener
		err error
	)

	err = c.Bind(&ct)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := UrlShotener.UrlShotenerUsecase.UpdateOne(ctx, &ct, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
