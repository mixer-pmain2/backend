package controller

import (
	"pmain2/internal/database"
	"pmain2/internal/models"
	"pmain2/pkg/cache"
	"time"
)

type spr struct{}

func initSprController() *spr {
	return &spr{}
}

func (m *spr) GetPodr() (*map[int]string, error) {
	cacheName := "spr_podr"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[int]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetPodr()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetPrava() (*[]models.PravaDict, error) {
	cacheName := "spr_prava"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]models.PravaDict)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetPrava()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}
