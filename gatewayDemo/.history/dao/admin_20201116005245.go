package dao

type AdminInfo struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"username" gorm:"column:username" description:"管理员用户名"`
	Salt 	  string	`json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
}

func ()  {
	
}