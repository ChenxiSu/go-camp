package service

import "hm1/model"
import "hm1/dao"

func GetUserById(id int64) (*model.User, error) {
	return dao.GetUserById(id)
}