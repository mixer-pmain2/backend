package controller

import (
	"fmt"
	"pmain2/internal/models"
	"pmain2/internal/types"
	"pmain2/pkg/cache"
	"strconv"
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
		ERROR.Println(err)
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (u *user) GetPrava(id int, isCache bool) (*map[int]int, error) {
	cacheName := fmt.Sprintf("user_prava_%v%h", id)

	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			res := item.(*map[int]int)
			return res, nil
		}
	}

	model := models.Model.User
	data, err := model.GetPrava(id)
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (u *user) ChangePassword(data types.ChangePassword) (int, error) {
	//data.Password, _ = utils.ToWin1251(data.Password)
	//data.NewPassword, _ = utils.ToWin1251(data.NewPassword)

	model := models.Model.User
	isAuth, err := model.UserAuth(strconv.FormatInt(data.UserId, 10), data.Password)
	if err != nil {
		ERROR.Println(err)
		return -1, err
	}
	if !isAuth {
		return 600, nil
	}

	_, err = model.ChangePassword(data)
	if err != nil {
		ERROR.Println(err)
		return -1, err
	}
	return 0, nil
}
