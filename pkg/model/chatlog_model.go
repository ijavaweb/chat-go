package model


type ChatLog struct {
	Id    int32  `gorm:"column:id" ,json:"id"`
	Who   string `gorm:"column:who" ,json:"who"`
	Query string `gorm:"column:query" ,json:"query"`
	Reply string `gorm:"column:reply" ,json:"reply"`
	Ctime int64  `gorm:"column:ctime" ,json:"ctime"`
}
