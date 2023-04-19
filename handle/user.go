package handle

import (
	"ggg/structs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (b *Backend) GetUserList(c *gin.Context) {

	var u_data []structs.User_List

	pageNo := c.Query("pageNo")
	pageSize := c.Query("pageSize")
	size, _ := strconv.Atoi(pageSize)
	no, _ := strconv.Atoi(pageNo)
	pageStart := (no - 1) * size

	var data_count int
	err := b.DbName.Get(&data_count, "Select count(1) from user")
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: u_data})
		return
	}
	err = b.DbName.Select(&u_data, "Select id,account,nickName,role,phone,password from user limit ?,?", pageStart, size)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: u_data})
		return
	}
	var records_data structs.Records_Data
	records_data.Total = data_count
	records_data.Page = no
	records_data.Data = u_data

	c.JSON(200, structs.Ret_Page{Status: 200, Message: "请求成功", Success: true, Records: records_data})
}

func (b *Backend) DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, structs.Ret{Status: 401, Message: "id不能为空", Success: false, Data: nil})
		return
	}
	_, err := b.DbName.Exec("Delete from user where id=?", id)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: nil})
		return
	}
	c.JSON(200, structs.Ret{Status: 200, Message: "删除成功", Success: true, Data: nil})
}

func (b *Backend) AddUser(c *gin.Context) {
	sqlStr := "insert into user(account,password,phone,nickName,role) values(?,?,?,?,?)"
	var body map[string]string
	_ = c.BindJSON(&body)

	acc := body["account"]
	pwd := body["password"]
	phone := body["phone"]
	nn := body["nickName"]
	role := body["role"]

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

	_, err := b.DbName.Exec(sqlStr, acc, pwd, phone, nn, role)
	if err != nil {
		c.JSON(500, structs.Ret{Status: 500, Message: "注册失败,系统错误", Success: true, Data: nil})
		return
	}
	c.JSON(200, structs.Ret{Status: 200, Message: "注册成功", Success: true, Data: nil})
}

func (b *Backend) UpdateUser(c *gin.Context) {
	sqlStr := "update user set account = ?, phone = ?,nickName=?,role=? where id=?"
	var body map[string]string
	_ = c.BindJSON(&body)

	id := body["id"]
	acc := body["account"]
	phone := body["phone"]
	nn := body["nickName"]
	role := body["role"]
	if id == "" {
		c.JSON(500, structs.Ret{Status: 401, Message: "id不能为空", Success: false, Data: nil})
		return
	}

	_, err := b.DbName.Exec(sqlStr, acc, phone, nn, role, id)
	if err != nil {
		c.JSON(500, structs.Ret{Status: 500, Message: "修改失败,系统错误", Success: false, Data: nil})
		return
	}
	c.JSON(200, structs.Ret{Status: 200, Message: "修改成功", Success: true, Data: nil})
}
