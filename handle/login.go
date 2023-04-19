package handle

import (
	"fmt"
	"ggg/structs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Backend struct {
	*structs.SqlInterface
}

func InitInterface(url string) (*Backend, error) {
	s, err := sqlx.Connect("mysql", url)
	if err != nil {
		return nil, err
	}
	return &Backend{&structs.SqlInterface{DbName: s, DbUrl: url}}, nil
}

func (b *Backend) UserLogin(c *gin.Context) {

	var u_data []structs.User_Login
	var body map[string]string
	_ = c.BindJSON(&body)

	acc := body["account"]
	pwd := body["password"]

	b.DbName.Select(&u_data, "Select account,nickName,role,phone from user where account=? and password=?", acc, pwd)
	fmt.Printf("data:%v", u_data)
	if len(u_data) < 1 {
		c.JSON(200, structs.Ret{Status: 500, Message: "账户名或密码错误", Success: false, Data: nil})
		return
	} else if len(u_data) > 1 {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误,查询到多个账户", Success: false, Data: nil})
		return
	}

	c.JSON(200, structs.Ret{Status: 200, Message: "登录成功", Success: true, Data: u_data[0]})
}

func (b *Backend) UserRegister(c *gin.Context) {
	sqlStr := "insert into user(account,password,phone,nickName) values(?,?,?,?)"
	var body map[string]string
	_ = c.BindJSON(&body)

	acc := body["account"]
	pwd := body["password"]
	phone := body["phone"]
	nn := body["nickName"]

	if acc == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "账号不能为空", Success: false, Data: nil})
		return
	} else if pwd == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "密码不能为空", Success: false, Data: nil})
		return
	} else if phone == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "手机号不能为空", Success: false, Data: nil})
		return
	} else if len(pwd) < 6 {
		c.JSON(200, structs.Ret{Status: 200, Message: "密码不能小于6位", Success: false, Data: nil})
		return
	} else if nn == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "用户昵称不能为空", Success: false, Data: nil})
		return
	}

	var u_data []structs.User
	var u_data_phone []structs.User
	var u_data_nn []structs.User
	b.DbName.Select(&u_data, "Select * from user where account=?", acc)
	b.DbName.Select(&u_data_phone, "Select * from user where phone=?", phone)
	b.DbName.Select(&u_data_nn, "Select * from user where nickName=?", nn)
	if len(u_data) > 0 {
		c.JSON(200, structs.Ret{Status: 200, Message: "账户名重复", Success: false, Data: nil})
		return
	} else if len(u_data_phone) > 0 {
		c.JSON(200, structs.Ret{Status: 200, Message: "手机号重复", Success: false, Data: nil})
		return
	} else if len(u_data_nn) > 0 {
		c.JSON(200, structs.Ret{Status: 200, Message: "昵称重复", Success: false, Data: nil})
		return
	}

	_, err := b.DbName.Exec(sqlStr, acc, pwd, phone, nn)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	c.JSON(200, structs.Ret{Status: 200, Message: "注册成功", Success: true, Data: nil})
}
