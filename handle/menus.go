package handle

import (
	"fmt"
	"ggg/structs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (b *Backend) GetMenuList(c *gin.Context) {

	var u_data []structs.Menu_Data

	pageNo := c.Query("pageNo")
	pageSize := c.Query("pageSize")
	size, _ := strconv.Atoi(pageSize)
	no, _ := strconv.Atoi(pageNo)
	pageStart := (no - 1) * size

	var data_count int
	err := b.DbName.Get(&data_count, "Select count(1) from menus")
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: err})
		return
	}
	err = b.DbName.Select(&u_data, "Select * from menus limit ?,?", pageStart, size)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: err})
		return
	}

	var records_data structs.Records_Data
	records_data.Total = data_count
	records_data.Page = no
	records_data.Data = u_data

	c.JSON(200, structs.Ret_Page{Status: 200, Message: "请求成功", Success: true, Records: records_data})
}

func (b *Backend) AddMenu(c *gin.Context) {
	var body map[string]string
	err := c.BindJSON(&body)

	if err != nil {
		c.JSON(200, structs.Ret{Status: 200, Message: "参数错误", Success: false, Data: nil})
		return
	}

	name := body["menuName"]
	key := body["menuKey"]
	icon := body["menuIcon"]
	sort := body["menuSort"]
	path := body["pagePath"]
	per := body["rolePermiss"]
	sta := body["status"]
	parent_menu := body["parentMenu"]

	if name == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单名称不能为空", Success: false, Data: nil})
		return
	} else if key == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单key不能为空", Success: false, Data: nil})
		return
	} else if path == "" {
		c.JSON(200, structs.Ret{Status: 200, Message: "页面路径不能为空", Success: false, Data: nil})
		return
	}

	var u_data_key int
	var u_data_name int
	var u_data_path int
	b.DbName.Get(&u_data_key, "Select count(*) from menus where menuKey=?", key)
	b.DbName.Get(&u_data_name, "Select count(*) from menus where menuName=?", name)
	b.DbName.Get(&u_data_path, "Select count(*) from menus where menuPath=?", path)
	if u_data_key > 0 {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单key重复", Success: false, Data: nil})
		return
	} else if u_data_name > 0 {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单名重复", Success: false, Data: nil})
		return
	} else if u_data_path > 0 {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单路径重复", Success: false, Data: nil})
		return
	}

	var u_data_parent []structs.Menu_Data
	b.DbName.Select(&u_data_parent, "Select *  from menus where menuKey=?", parent_menu)

	if parent_menu != "layout" && u_data_parent == nil {
		c.JSON(200, structs.Ret{Status: 200, Message: "未找到父级目录", Success: false, Data: nil})
		return
	}

	sqlStr2 := "update menus set children = ? where menuKey=?"

	var new_children string

	if u_data_parent == nil {
		new_children = key
	} else {
		new_children = u_data_parent[0].Children + "," + key
	}

	_, err = b.DbName.Exec(sqlStr2, new_children, parent_menu)

	if err != nil {
		c.JSON(500, structs.Ret{Status: 200, Message: "更新children失败", Success: true, Data: nil})
		return
	}
	sqlStr := "insert into menus(menuName,menuKey,menuIcon,pagePath,rolePermiss,menuSort,status,children) values(?,?,?,?,?,?,?,?)"
	_, err = b.DbName.Exec(sqlStr, name, key, icon, path, per, sort, sta, "")
	if err != nil {
		c.JSON(500, structs.Ret{Status: 200, Message: "添加菜单失败", Success: true, Data: nil})
		return
	}

	c.JSON(200, structs.Ret{Status: 200, Message: "添加成功", Success: true, Data: nil})
}

func (b *Backend) UpdateMenu(c *gin.Context) {
	var body map[string]string
	err := c.BindJSON(&body)
	var u_data []structs.Menu_Data
	if err != nil {
		c.JSON(200, structs.Ret{Status: 200, Message: "参数错误", Success: false, Data: nil})
		return
	}

	id := body["id"]
	name := body["menuName"]
	key := body["menuKey"]
	icon := body["menuIcon"]
	sort := body["menuSort"]
	path := body["pagePath"]
	per := body["rolePermiss"]
	sta := body["status"]

	err = b.DbName.Select(&u_data, "select * from menus where id=?", id)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: nil})
		return
	}

	if u_data[0].Status == "2" {
		c.JSON(200, structs.Ret{Status: 200, Message: "无法更改锁定数据", Success: false, Data: nil})
		return
	}

	var u_data_key int
	var u_data_name int
	var u_data_path int
	b.DbName.Get(&u_data_key, "Select count(*) from menus where menuKey=?", key)
	b.DbName.Get(&u_data_name, "Select count(*) from menus where menuName=?", name)
	b.DbName.Get(&u_data_path, "Select count(*) from menus where menuPath=?", path)
	if u_data_key > 1 {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单key重复", Success: false, Data: nil})
		return
	} else if u_data_name > 1 {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单名重复", Success: false, Data: nil})
		return
	} else if u_data_path > 1 {
		c.JSON(200, structs.Ret{Status: 200, Message: "菜单路径重复", Success: false, Data: nil})
		return
	}

	_, err = b.DbName.Exec("update menus set menuName=?,menuKey=?,menuIcon=?,pagePath=?,menuSort=?,rolePermiss=?,status=? where id=?", name, key, icon, path, sort, per, sta, id)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "编辑错误", Success: false, Data: nil})
		return
	}
	c.JSON(200, structs.Ret{Status: 200, Message: "修改成功", Success: true, Data: nil})
}

func (b *Backend) DeleteMenu(c *gin.Context) {
	var u_data []structs.Menu_Data
	id := c.Query("id")
	err := b.DbName.Select(&u_data, "Select * from menus where id = ?", id)
	fmt.Printf("%v:::", &u_data)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: err})
		return
	}
	if u_data[0].Status == "2" {
		c.JSON(200, structs.Ret{Status: 500, Message: "无法删除锁定数据", Success: false, Data: err})
		return
	}
	_, err = b.DbName.Exec("Delete from menus where id = ?", id)
	if err != nil {
		c.JSON(200, structs.Ret{Status: 500, Message: "系统错误", Success: false, Data: err})
		return
	}
	c.JSON(200, structs.Ret{Status: 200, Message: "删除成功", Success: true, Data: nil})
}
