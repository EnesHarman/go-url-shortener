package repository

import (
	"com.github/EnesHarman/url-shortener/internal/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type UrlRepository interface {
	AddNewUrl(url *model.Url) error
	GetUrlByShortUrl(shortUrl string) (*model.Url, error)
	GetUrls(page int, size int) ([]model.Url, error)
	DeleteUrl(id int) error
}

type UrlRepositoryImpl struct {
	dbPool *pgxpool.Pool
}

func NewUrlRepository(dbPool *pgxpool.Pool) UrlRepository {
	return &UrlRepositoryImpl{
		dbPool: dbPool,
	}
}

func (repo UrlRepositoryImpl) GetUrlByShortUrl(shortUrl string) (*model.Url, error) {
	ctx := context.Background()
	row := repo.dbPool.QueryRow(ctx, "SELECT realUrl, shortUrl, expiredate, ts FROM url WHERE shortUrl = $1", shortUrl)
	url := model.Url{}
	err := row.Scan(&url.RealUrl, &url.ShortUrl, &url.ExpireDate, &url.Ts)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (repo UrlRepositoryImpl) GetUrls(page int, size int) ([]model.Url, error) {
	var urls []model.Url
	ctx := context.Background()
	row, err := repo.dbPool.Query(ctx, "SELECT id ,realUrl, shortUrl, expiredate, ts FROM url LIMIT $1 OFFSET $2", size, (page-1)*size)
	if err != nil {
		log.Error("Failed to get urls: %v", err)
		return nil, errors.New("Failed to get urls")
	}
	for row.Next() {
		url := model.Url{}
		err := row.Scan(&url.Id, &url.RealUrl, &url.ShortUrl, &url.ExpireDate, &url.Ts)
		if err != nil {
			log.Error("Failed to scan url: %v", err)
			return nil, errors.New("Failed to get urls")
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (repo UrlRepositoryImpl) AddNewUrl(url *model.Url) error {
	ctx := context.Background()
	_, err := repo.dbPool.Exec(ctx, "INSERT INTO url (realUrl, shortUrl, expiredate, ts) VALUES ($1, $2, $3, $4)", url.RealUrl, url.ShortUrl, url.ExpireDate, url.Ts)
	if err != nil {
		log.Error("Failed to insert new url: %v", err)
		return err
	}
	return nil
}

func (repo UrlRepositoryImpl) DeleteUrl(id int) error {
	ctx := context.Background()
	res, err := repo.dbPool.Exec(ctx, "DELETE FROM url WHERE id = $1", id)
	if err != nil {
		log.Error("Failed to delete url: %v", err)
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("Url not found")
	}
	return nil
}
