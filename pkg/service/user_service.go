package service

import (
	"blog-go/pkg/db"
	"blog-go/pkg/model"
	"log"
)

func CreateUser(user *model.UserInfo) error {
	user.Status=1
	err:=db.CreateUser(user)
	return err
}
func UpdateUser(user *model.UserInfo) error  {
	err:=db.UpdateUser(user)
	return err
}
func DeleteUser(user *model.UserInfo) error  {
	user.Status=0
	err:=db.DeleteUser(user)
	return err
}
func CreateWechatUser(username string)  {
	u, err := db.GetUser(username)
	if err != nil {
		return
	}
	if u != nil {
		return
	}
	user := &model.UserInfo{
		Id:       0,
		UserName: username,
		Password: "",
		Role:     "wechat",
		Status:   0,
	}
	err = db.CreateUser(user)
	if err != nil {
		log.Println(err.Error())
	}
}