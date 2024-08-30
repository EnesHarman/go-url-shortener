package controllers

import (
	"com.github/EnesHarman/url-shortener/internal/model"
	"com.github/EnesHarman/url-shortener/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UrlController interface {
	AddUrl(c *gin.Context)
	RegisterRoutes(e *gin.Engine)
	RedirectUrl(c *gin.Context)
}

type UrlControllerImpl struct {
	urlService services.UrlService
}

func NewUrlController(urlService services.UrlService) UrlController {
	return &UrlControllerImpl{
		urlService: urlService,
	}
}

func (ctrl UrlControllerImpl) RegisterRoutes(e *gin.Engine) {
	e.POST("/url/add", ctrl.AddUrl)
	e.GET("/:shortUrl", ctrl.RedirectUrl)
	e.GET("/url/list", ctrl.ListUrls)
	e.DELETE("/url/delete/:id", ctrl.DeleteUrl)
}

func (ctrl UrlControllerImpl) AddUrl(c *gin.Context) {
	var url model.Url
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ctrl.urlService.AddUrl(&url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shortUrl": res})
}

func (ctrl UrlControllerImpl) ListUrls(c *gin.Context) {
	page, size := 1, 10
	if pg := c.Query("page"); pg != "" {
		if parsedPage, err := strconv.Atoi(pg); err == nil {
			page = parsedPage
		}
	}

	if ps := c.Query("size"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil {
			size = parsedPageSize
		}
	}
	urls, err := ctrl.urlService.GetUrls(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": urls})
}

func (ctrl UrlControllerImpl) DeleteUrl(c *gin.Context) {
	var id int
	idStr := c.Params.ByName("id")
	if parsedPage, err := strconv.Atoi(idStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	} else {
		id = parsedPage
	}
	err := ctrl.urlService.DeleteUrl(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Url deleted successfully"})
}

func (ctrl UrlControllerImpl) RedirectUrl(c *gin.Context) {
	shortUrl := c.Params.ByName("shortUrl")
	url, err := ctrl.urlService.GetUrlByShortUrl(shortUrl)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url.RealUrl)
}
