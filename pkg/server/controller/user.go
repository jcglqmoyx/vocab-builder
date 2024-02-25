package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"vocab-builder/pkg/server/conf"
	"vocab-builder/pkg/server/dao"
	"vocab-builder/pkg/server/model"
	"vocab-builder/pkg/server/util"
)

func Register(c *gin.Context) {
	type RegisterForm struct {
		Username        string `form:"username" json:"username"`
		Email           string `form:"email" json:"email"`
		Avatar          string `form:"avatar" json:"avatar"`
		Password        string `form:"password" json:"password"`
		ConfirmPassword string `form:"confirm_password" json:"confirm_password"`
	}
	var registerForm RegisterForm
	if err := c.ShouldBind(&registerForm); err != nil {
		util.JsonHttpResponse(c, 1, "参数错误", nil)
		return
	}
	if registerForm.Username == "" || registerForm.Password == "" || registerForm.ConfirmPassword == "" {
		util.JsonHttpResponse(c, 1, "注册失败, 用户名或密码不能为空", nil)
		return
	}
	if registerForm.ConfirmPassword != registerForm.Password {
		util.JsonHttpResponse(c, 1, "注册失败, 两次输入的密码不一致", nil)
		return
	}
	if _, found := dao.FindUserByUsername(registerForm.Username); found {
		util.JsonHttpResponse(c, 1, "该用户名已被占用", nil)
		return
	}
	if _, found := dao.FindUserByEmail(registerForm.Email); found {
		util.JsonHttpResponse(c, 1, "该邮箱已被占用", nil)
		return
	}
	user := model.User{
		Username: registerForm.Username,
		Email:    registerForm.Email,
		Avatar:   registerForm.Avatar,
		Password: registerForm.Password,
	}
	defaultDictionaries := conf.Cfg.Dictionary.Dictionaries
	dao.AddUser(&user)
	for _, dictionary := range defaultDictionaries {
		dao.AddDictionary(&model.Dictionary{
			UserID: user.ID,
			Title:  dictionary.Title,
			Prefix: dictionary.Prefix,
			Suffix: dictionary.Suffix,
		})
	}
	util.JsonHttpResponse(c, 0, "注册成功", nil)
}

func Login(c *gin.Context) {
	var dto *model.User
	if err := c.ShouldBind(&dto); err != nil {
		util.JsonHttpResponse(c, 1, "参数错误", nil)
		return
	}
	user, found := dao.FindUserByUsername(dto.Username)
	if !found {
		util.JsonHttpResponse(c, 1, "用户名不存在", nil)
		return
	}
	if util.HashPassword(dto.Password, user.PasswordSalt) != user.Password {
		util.JsonHttpResponse(c, 1, "密码错误", nil)
		return
	}
	token, _ := util.GenerateJWT(user.ID)
	json := map[string]string{
		"token":   token,
		"user_id": strconv.Itoa(user.ID),
	}
	util.JsonHttpResponse(c, 0, "登录成功", json)
}

func DeleteUserByID(c *gin.Context) {
	if userID, err := strconv.Atoi(c.Param("id")); err != nil {
		util.JsonHttpResponse(c, 1, "id值不合法", nil)
	} else {
		if _, found := dao.FindUserByID(userID); !found {
			util.JsonHttpResponse(c, 1, "该用户不存在, 无法删除", nil)
		} else {
			dao.DeleteUserByID(userID)
			util.JsonHttpResponse(c, 0, "用户删除成功", nil)
		}
	}
}

func UpdateUser(c *gin.Context) {
	userID, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 2, err.Error(), nil)
		return
	} else {
		if _, found := dao.FindUserByID(userID); !found {
			util.JsonHttpResponse(c, 1, "用户不存在, 无法更新", nil)
			return
		}
		var user *model.User
		if err := c.ShouldBind(&user); err != nil {
			util.JsonHttpResponse(c, 1, "参数不合法", nil)
		} else {
			book, found := dao.FindBookByID(user.CurrentBookID)
			if !found {
				util.JsonHttpResponse(c, 1, "词书不存在", nil)
				return
			}
			if book.UserID != userID {
				util.JsonHttpResponse(c, 1, "词书不属于该用户", nil)
				return
			}

			if user.DailyCount <= 0 {
				util.JsonHttpResponse(c, 1, "每日学习量不合法", nil)
				return
			}

			if user.DailyCount > 500000 {
				util.JsonHttpResponse(c, 1, "每日学习量过大", nil)
				return
			}

			if user.ReviewFrequencyFormula == "" {
				util.JsonHttpResponse(c, 1, "复习频率公式不能为空", nil)
				return
			}

			if user.TimesCountedAsKnown <= 0 {
				util.JsonHttpResponse(c, 1, "\"点击几次“认识”算学会\"参数不合法", nil)
				return
			}

			if user.TimesCountedAsKnown > 1000 {
				util.JsonHttpResponse(c, 1, "\"点击几次“认识”算学会\"参数不合法", nil)
				return
			}

			_, ok := util.ParseReviewFrequencyFormula(user.ReviewFrequencyFormula)
			if !ok {
				util.JsonHttpResponse(c, 1, "复习频率公式不合法", nil)
				return
			}
			user.ID = userID
			dao.UpdateUser(user)
			util.JsonHttpResponse(c, 0, "success", user)
		}
	}
}

func GetUserProfile(c *gin.Context) {
	id, err := util.GetUserID(c)
	if err != nil {
		util.JsonHttpResponse(c, 1, err.Error(), nil)
		return
	}
	user, found := dao.FindUserByID(id)
	if found {
		userInfo := model.User{
			Username:               user.Username,
			Email:                  user.Email,
			Avatar:                 user.Avatar,
			CurrentBookID:          user.CurrentBookID,
			DailyCount:             user.DailyCount,
			TimesCountedAsKnown:    user.TimesCountedAsKnown,
			ReviewFrequencyFormula: user.ReviewFrequencyFormula,
		}
		util.JsonHttpResponse(c, 0, "success", userInfo)
	} else {
		util.JsonHttpResponse(c, 1, "用户不存在", nil)
	}
}
