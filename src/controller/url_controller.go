package controller

import (
	"github.com/gin-gonic/gin"
	"url-shortener/src/models/urls"
)

type imageController struct {
	service urls.Service
}

// @BasePath /url-shortener

// RegisterImage @Description get data by Image rul
// @Accept  json
// @Produce  json
// @Param   image_url     query    string     true    "https://google.com/img.jpeg"
// @Success 200 {string} json	"ok"
// @Failure 400 {object} rest_error.RestErr "We need image_url!!"
// @Failure 404 {object} rest_error.RestErr "Can not find image_url"
// @Router /register_image [get]
func (ic imageController) RegisterImage(c *gin.Context) {
	args := struct {
		ImageUrl string `form:"image_url" binding:"required"`
	}{}

	if err := c.BindQuery(&args); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error()})
		return
	}

	//r,err := ic.service.GetUrl(args.ImageUrl)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Message()})
	//	return
	//}
	//c.JSON(http.StatusOK, r)
}

// @BasePath /url-shortener

// GetUrl @Description get data by Image rul
// @Accept  json
// @Produce  json
// @Param   image_url     query    string     true    "https://google.com/img.jpeg"
// @Success 200 {string} data	"ok"
// @Failure 400 {object} rest_error.RestErr "We need image_url!!"
// @Failure 404 {object} rest_error.RestErr "Can not find image_url"
// @Router /get_image [get]
func (ic imageController) GetUrl(c *gin.Context) {
	args := struct {
		ImageUrl string `form:"image_url" binding:"required"`
	}{}

	if err := c.BindQuery(&args); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error()})
		return
	}

	//d,err:= ic.service.GetUrl(args.ImageUrl)
	//
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error":  fmt.Sprintf("Unable to decode JSON request body: %v", err)})
	//	return
	//}

	//urls, err := ic.service.GetUrl(args.ImageUrl)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Message()})
	//	return
	//}
	//c.JSON(http.StatusOK, string(d))
}

type SiteMetaInterfaces interface {
	GetUrl(c *gin.Context)
	RegisterImage(c *gin.Context)
	Route(prefix *gin.RouterGroup)
}

func NewHandler(service urls.Service) SiteMetaInterfaces {
	return &imageController{service: service}
}

func (ic *imageController) Route(prefix *gin.RouterGroup) {
	prefix.GET(
		"/get_url",
		ic.GetUrl,
	)
	prefix.GET(
		"/register_url",
		ic.RegisterImage,
	)
}
