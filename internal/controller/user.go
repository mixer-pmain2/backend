package controller

import (
	"fmt"
	"pmain2/internal/models"
	"pmain2/pkg/cache"
)

type user struct{}

func initUserController() *user {
	return &user{}
}

func (u *user) IsAuth(login, password string) (bool, error) {
	model := models.Model.User
	ok, err := model.UserAuth(login, password)
	if err != nil {
		ERROR.Println(err.Error())
		return false, err
	}
	return ok, nil
}

func (u *user) GetUch(id int) (*map[int][]int, error) {
	cacheName := fmt.Sprintf("user_%v%_uch", id)

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[int][]int)
		return res, nil
	}

	model := models.Model.User
	data, err := model.GetUch(id)
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}
