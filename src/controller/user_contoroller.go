package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"url-shortener/config"
	"url-shortener/src/models/url_shortener"
	"url-shortener/src/models/users"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	userJwt.GET("/users", handler.GetAll)
	userJwt.PUT("/user", handler.UpdateOne)
}

func isRequestValid(m *users.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (user *UserHandler) Redirect(c echo.Context) error {
	url := c.Param("key")
	ctx := c.Request().Context()
	OriginalURL, err := user.UrlUsecase.FindOneByKey(ctx, url)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.Redirect(http.StatusMovedPermanently, OriginalURL)
}
func (user *UserHandler) InsertOne(c echo.Context) error {
	var (
		usr users.User
		err error
	)

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

	result, err := user.UsrUsecase.InsertOne(ctx, &usr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

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

func (user *UserHandler) GetAll(c echo.Context) error {

	type Response struct {
		Total       int64        `json:"total"`
		PerPage     int64        `json:"per_page"`
		CurrentPage int64        `json:"current_page"`
		LastPage    int64        `json:"last_page"`
		From        int64        `json:"from"`
		To          int64        `json:"to"`
		User        []users.User `json:"users"`
	}

	var (
		res   []users.User
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

	res, count, err = user.UsrUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
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
		User:        res,
	}

	return c.JSON(http.StatusOK, result)
}

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
