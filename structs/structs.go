package structs

import "github.com/jmoiron/sqlx"

type User_List struct {
	Id       string `db:"id" json:"id"`
	NickName string `db:"nickName" json:"nickName"`
	Phone    string `db:"phone" json:"phone"`
	Account  string `db:"account" json:"account"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type SqlInterface struct {
	DbName *sqlx.DB
	DbUrl  string
}

type User struct {
	Id       int    `db:"id" json:"id"`
	NickName string `db:"nickName" json:"nickName"`
	Phone    string `db:"phone" json:"phone"`
	Account  string `db:"account" json:"account"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}
type User_Login struct {
	NickName string `db:"nickName" json:"nickName"`
	Phone    string `db:"phone" json:"phone"`
	Account  string `db:"account" json:"account"`
	Role     string `db:"role" json:"role"`
}

type Ret struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}
type Ret_Page struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Records any    `json:"records"`
}

type Records_Data struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Data  any `json:"data"`
}

type Menu_Data struct {
	Id          string `db:"id" json:"id"`
	MenuName    string `db:"menuName" json:"menuName"`
	MenuKey     string `db:"menuKey" json:"menuKey"`
	MenuIcon    string `db:"menuIcon" json:"menuIcon"`
	PagePath    string `db:"pagePath" json:"pagePath"`
	MenuSort    string `db:"menuSort" json:"menuSort"`
	RolePermiss string `db:"rolePermiss" json:"rolePermiss"`
	Status      string `db:"status" json:"status"`
	Children    string `db:"children" json:"children"`
}
