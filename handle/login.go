package handle

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

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

func InitInterface(url string) (*SqlInterface, error) {
	s, err := sqlx.Connect("mysql", url)
	if err != nil {
		return nil, err
	}
	return &SqlInterface{DbName: s, DbUrl: url}, nil
}

func (si *SqlInterface) UserLogin(c *gin.Context) {

	var u_data []User_Login
	var body map[string]string
	_ = c.BindJSON(&body)

	acc := body["account"]
	pwd := body["password"]

	si.DbName.Select(&u_data, "Select account,nickName,role,phone from user where account=? and password=?", acc, pwd)
	fmt.Printf("data:%v", u_data)
	if len(u_data) < 1 {
		c.JSON(200, Ret{Status: 500, Message: "账户名或密码错误", Success: false, Data: nil})
		return
	} else if len(u_data) > 1 {
		c.JSON(200, Ret{Status: 500, Message: "系统错误,查询到多个账户", Success: false, Data: nil})
		return
	}

	c.JSON(200, Ret{Status: 200, Message: "登录成功", Success: true, Data: u_data[0]})
}

func (si *SqlInterface) UserRegister(c *gin.Context) {
	sqlStr := "insert into user(account,password,phone,nickName) values(?,?,?,?)"
	var body map[string]string
	_ = c.BindJSON(&body)

	acc := body["account"]
	pwd := body["password"]
	phone := body["phone"]
	nn := body["nickName"]

	if acc == "" {
		c.JSON(200, Ret{Status: 200, Message: "账号不能为空", Success: false, Data: nil})
		return
	} else if pwd == "" {
		c.JSON(200, Ret{Status: 200, Message: "密码不能为空", Success: false, Data: nil})
		return
	} else if phone == "" {
		c.JSON(200, Ret{Status: 200, Message: "手机号不能为空", Success: false, Data: nil})
		return
	} else if len(pwd) < 6 {
		c.JSON(200, Ret{Status: 200, Message: "密码不能小于6位", Success: false, Data: nil})
		return
	} else if nn == "" {
		c.JSON(200, Ret{Status: 200, Message: "用户昵称不能为空", Success: false, Data: nil})
		return
	}

	var u_data []User
	var u_data_phone []User
	var u_data_nn []User
	si.DbName.Select(&u_data, "Select * from user where account=?", acc)
	si.DbName.Select(&u_data_phone, "Select * from user where phone=?", phone)
	si.DbName.Select(&u_data_nn, "Select * from user where nickName=?", nn)
	if len(u_data) > 0 {
		c.JSON(200, Ret{Status: 200, Message: "账户名重复", Success: false, Data: nil})
		return
	} else if len(u_data_phone) > 0 {
		c.JSON(200, Ret{Status: 200, Message: "手机号重复", Success: false, Data: nil})
		return
	} else if len(u_data_nn) > 0 {
		c.JSON(200, Ret{Status: 200, Message: "昵称重复", Success: false, Data: nil})
		return
	}

	_, err := si.DbName.Exec(sqlStr, acc, pwd, phone, nn)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	c.JSON(200, Ret{Status: 200, Message: "注册成功", Success: true, Data: nil})

}
