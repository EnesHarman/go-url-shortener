package services

import (
	"com.github/EnesHarman/url-shortener/config"
	"com.github/EnesHarman/url-shortener/internal/kafka"
	kafkaModel "com.github/EnesHarman/url-shortener/internal/kafka/model"
	"com.github/EnesHarman/url-shortener/internal/model"
	"com.github/EnesHarman/url-shortener/internal/repository"
	"fmt"
	"math/rand"
	"time"
)

type UrlService interface {
	AddUrl(url *model.Url) (string, error)
	GetUrlByShortUrl(shortUrl string) (*model.Url, error)
	GetUrls(page int, size int) ([]model.Url, error)
	DeleteUrl(id int) error
	PublishClickEvent(urlId int, userId string)
}

type UrlServiceImpl struct {
	repository    repository.UrlRepository
	urlConfig     config.UrlShortenerConfig
	eventProducer kafka.ClickEventProducer
}

func NewUrlService(repository repository.UrlRepository, urlConfig config.UrlShortenerConfig, eventProducer kafka.ClickEventProducer) UrlService {
	return &UrlServiceImpl{
		repository:    repository,
		urlConfig:     urlConfig,
		eventProducer: eventProducer,
	}
}

func (service UrlServiceImpl) GetUrlByShortUrl(shortUrl string) (*model.Url, error) {
	return service.repository.GetUrlByShortUrl(shortUrl)
}

func (service UrlServiceImpl) GetUrls(page int, size int) ([]model.Url, error) {
	return service.repository.GetUrls(page, size)
}

func (service UrlServiceImpl) AddUrl(url *model.Url) (string, error) {
	url.Ts = time.Now()
	url.ShortUrl = service.generateRandomUrl()
	if err := service.repository.AddNewUrl(url); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", service.urlConfig.BaseUrl, url.ShortUrl), nil

}

func (service UrlServiceImpl) DeleteUrl(id int) error {
	return service.repository.DeleteUrl(id)
}

func (service UrlServiceImpl) PublishClickEvent(urlId int, userId string) {
	event := kafkaModel.ClickEvent{
		UserId: userId,
		Ts:     time.Now(),
		UrlId:  urlId,
	}
	service.eventProducer.Produce(event)
}

func (service UrlServiceImpl) generateRandomUrl() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, service.urlConfig.Length)

	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
