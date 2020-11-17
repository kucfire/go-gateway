package dto

import "time"

type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avator       string    `json:"avator"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}
