package db

import "blog-go/pkg/model"

func CreateChatLog(chatLog *model.ChatLog)  error{
	db:=DB.Table("chat_log").Create(chatLog)
	return db.Error
}