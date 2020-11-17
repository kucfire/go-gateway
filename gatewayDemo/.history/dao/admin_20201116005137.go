package dao

type AdminInfo struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	Username  string    `json:"username" gorm:"column:username" description:"管理员用户名"`
	Salt int	`json:"salt" gorm:"column:salt" description:"盐"`
	Password    int       `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
}

func ()  {
	
}