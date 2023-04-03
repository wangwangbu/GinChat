package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// GetUserList godoc
// @Summary 所有用户
// @Tags         用户
// @Success      200  {string}  json{"code", "message"}
// @Router       /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(200, gin.H{
		"message": data,
	})
}

// CreateUser godoc
// @Summary 新增用户
// @Tags         用户
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success      200  {string}  json{"code", "message"}
// @Router       /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.NewUserBasic()
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())
	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
		})
		return
	}

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code": -1,
			"message": "该用户名已被占用",
		})
		return
	}
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code": 0,
		"message": "新增用户成功",
		"data": user,
	})
}

// FindUserByNameAndPwd godoc
// @Summary 登录验证
// @Tags         用户
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success      200  {string}  json{"code", "message"}
// @Router       /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"message": "该用户不存在",
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code": -1,
			"message": "密码不正确",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"message": "登录成功",
		"data": user,
	})
}

// CreateUser godoc
// @Summary 删除用户
// @Tags         用户
// @param id query string false "id"
// @Success      200  {string}  json{"code", "message"}
// @Router       /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.NewUserBasic()
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code": 0,
		"message": "删除用户成功",
	})
}

// CreateUser godoc
// @Summary 修改用户
// @Tags         用户
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "手机号"
// @param email formData string false "邮箱"
// @Success      200  {string}  json{"code", "message"}
// @Router       /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.NewUserBasic()
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": -1,
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code": 0,
			"message": "修改用户成功",
		})
	}
}
