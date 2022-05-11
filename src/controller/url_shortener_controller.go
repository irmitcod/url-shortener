package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	validator "gopkg.in/go-playground/validator.v9"
	"math"
	"net/http"
	"strconv"
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
	userJwt.GET("/UrlShotener", handler.FindOne)
	userJwt.GET("/UrlShoteners", handler.GetAll)
	userJwt.PUT("/UrlShotener", handler.UpdateOne)
}

func isRequestUrlValid(m *url_shortener.UrlShortener) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (UrlShotener *UrlShotenerHandler) InsertOne(c echo.Context) error {
	var (
		ct  url_shortener.UrlShortener
		err error
	)
	user := c.Get("user")
	token := user.(*users.User)

	err = c.Bind(&ct)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestUrlValid(&ct); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ct.UserID = token.ID
	result, err := UrlShotener.UrlShotenerUsecase.InsertOne(ctx, &ct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	result.ShortUrl = fmt.Sprintf("%s%s", UrlShotener.configuration.BaseUrl, result.ShortUrl)

	return c.JSON(http.StatusOK, result)
}

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

func (UrlShotener *UrlShotenerHandler) GetAll(c echo.Context) error {

	type Response struct {
		Total       int64                        `json:"total"`
		PerPage     int64                        `json:"per_page"`
		CurrentPage int64                        `json:"current_page"`
		LastPage    int64                        `json:"last_page"`
		From        int64                        `json:"from"`
		To          int64                        `json:"to"`
		UrlShotener []url_shortener.UrlShortener `json:"UrlShoteners"`
	}

	var (
		res   []url_shortener.UrlShortener
		count int64
	)

	rp, err := strconv.ParseInt(c.QueryParam("rp"), 10, 64)
	if err != nil {
		rp = 25
	}

	page, err := strconv.ParseInt(c.QueryParam("p"), 10, 64)
	if err != nil {
		page = 1
	}

	filters := bson.D{{"name", primitive.Regex{Pattern: ".*" + c.QueryParam("name") + ".*", Options: "i"}}}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, count, err = UrlShotener.UrlShotenerUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	result := Response{
		Total:       count,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(count) / float64(rp))),
		From:        page*rp - rp + 1,
		To:          page * rp,
		UrlShotener: res,
	}

	return c.JSON(http.StatusOK, result)
}

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
