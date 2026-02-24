package handlers

import (
	"net/http"
	"strconv"

	"PLMS/internal/models"
	"PLMS/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db      *gorm.DB
	service *services.UserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db:      db,
		service: services.NewUserService(db),
	}
}

// GetCurrentUser 获取当前登录用户信息（兼容性接口，与 /auth/current-user 返回相同格式）
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}

	authService := services.NewAuthService(h.db)
	user, err := authService.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
			"data":    nil,
		})
		return
	}

	// 返回与 Ant Design Pro 兼容的格式
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":              user.Name,
			"avatar":            user.ToUserInfo()["avatar"],
			"userid":            user.Username,
			"email":             "",
			"signature":         "",
			"title":             "",
			"group":             "",
			"tags":              []gin.H{},
			"notifyCount":       0,
			"unreadCount":       0,
			"country":           "China",
			"access":            user.Role,
			"isDefaultPassword": user.IsDefaultPassword == 1,
			"geographic": gin.H{
				"province": gin.H{"label": "浙江省", "key": "330000"},
				"city":     gin.H{"label": "杭州市", "key": "330100"},
			},
			"address": "",
			"phone":   "",
		},
	})
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// GetUser 获取单个用户
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "用户不存在",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	user.ID = uint(id)
	if err := h.service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}
