package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CheemsGoUp/Simplified-Douyin-Project/global"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = make(map[string]User)

// {
// 	"zhangleidouyin": {
// 		Id:            1,
// 		Name:          "zhanglei",
// 		FollowCount:   10,
// 		FollowerCount: 5,
// 		// IsFollow:      true,
// 	},
// }

// var userIdSequence = int64(1)

type UserRegisterRequest struct {
	Username string `json:"username"` // 注册用户名，最长32个字符
	Password string `json:"password"` // 密码，最长32个字符
}

type UserRegisterResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserId     int64  `json:"user_id"`     // 用户id
	Token      string `json:"token"`       // 用户鉴权token
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	// 接收注册信息
	username, ok1 := c.GetQuery("username")
	password, ok2 := c.GetQuery("password")
	if !ok1 || !ok2 {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  "Register: Cannot query username or password",
		})
		return
	}
	token := username + "_" + password

	// 检验注册信息合法性，用户名1-32个字符（不能含有"_"，确保token唯一），密码6-32个字符
	if len([]rune(username)) == 0 || len([]rune(username)) > 32 {
		go func() {
			fmt.Printf("%v\n", username)
		}()
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: 2,
			StatusMsg:  "Register: Username should be 1-32 characters",
		})
		return
	}
	if strings.ContainsAny(username, "_") {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: 3,
			StatusMsg:  "Register: Username should not contain '_'",
		})
		return
	}
	if len([]rune(password)) < 6 || len([]rune(password)) > 32 {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: 4,
			StatusMsg:  "Register: Password should be 6-32 characters",
		})
		return
	}

	// migrate到数据库
	global.DB.AutoMigrate(&User{})

	// 验证用户名是否已存在
	newUser := User{
		Name: username,
	}
	if !global.DB.NewRecord(newUser) {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: 5,
			StatusMsg:  "Register: User already exist",
		})
		return
	}

	// 如果为新用户，把用户数据存到数据库，并返回用户的Id
	newUser = User{
		Name:     username,
		Password: password,
	}
	err := global.DB.Create(&newUser).Error
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: 6,
			StatusMsg:  "Register: Fail to create new user in database",
		})
		return
	}
	// err = global.DB.Select("id").Where("name = ?", username).Last(&newUser).Error
	// if err != nil {
	// 	c.JSON(http.StatusOK, UserRegisterResponse{
	// 		StatusCode: 7,
	// 		StatusMsg:  "Fail to get user id",
	// 	})
	// 	return
	// }

	// 将用户信息登录信息存到 usersLoginInfo (map[token]User)
	usersLoginInfo[token] = newUser

	// 返回用户注册信息和token
	c.JSON(http.StatusOK, UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "Register: Success",
		UserId:     newUser.Id,
		Token:      token,
	})

	// if _, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	// 	})
	// } else {
	// 	atomic.AddInt64(&userIdSequence, 1)
	// 	newUser := User{
	// 		Id:   userIdSequence,
	// 		Name: username,
	// 	}
	// 	usersLoginInfo[token] = newUser
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   userIdSequence,
	// 		Token:    username + password,
	// 	})
	// }
}

func Login(c *gin.Context) {
	username, ok1 := c.GetQuery("username")
	enteredPassword, ok2 := c.GetQuery("password")
	if !ok1 || !ok2 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1,
				StatusMsg: "Login: Cannot query username or password"},
		})
		return
	}

	// migrate到数据库
	global.DB.AutoMigrate(&User{})

	// 根据输入的username从数据库调取user信息
	var user User
	err := global.DB.Where("name = ?", username).Last(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2,
				StatusMsg: "Login: Cannot find user"},
		})
		return
	}

	// 比对输入的密码和user信息中的密码是否一致
	if enteredPassword != user.Password {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 3,
				StatusMsg: "Login: Password is not correct"},
		})
		return
	}

	// 如果密码一致，则把user信息和token存放到usersLoginInfo（map[token]User)
	token := username + "_" + user.Password
	usersLoginInfo[token] = user
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Login: Success"},
		UserId:   user.Id,
		Token:    token,
	})

	// if user, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   user.Id,
	// 		Token:    token,
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	// 	})
	// }
}

func UserInfo(c *gin.Context) {
	// 获取token
	token, ok := c.GetQuery("token")
	if !ok {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "UserInfo: Cannot query token"},
		})
		return
	}

	// 根据token拉取用户信息
	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 2, StatusMsg: "UserInfo: User doesn't exist"},
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0, StatusMsg: "UserInfo: Success"},
		User:     user,
	})
}
