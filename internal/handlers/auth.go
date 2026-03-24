package handlers

import (
	"net/http"
	"strconv"

	"PLMS/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	db      *gorm.DB
	service *services.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		db:      db,
		service: services.NewAuthService(db),
	}
}

// Login 用户登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    result,
	})
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/auth/current-user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}

	user, err := h.service.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user.ToUserInfo(),
	})
}

// ChangePassword 修改密码
// POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}

	var req services.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	if err := h.service.ChangePassword(userID.(int64), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
		"data":    nil,
	})
}

// Logout 用户登出
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT是无状态的，登出主要在客户端清除token
	// 如果需要实现token黑名单，可以在这里将token加入黑名单
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
		"data":    nil,
	})
}

// CheckDefaultPassword 检查是否使用默认密码
// GET /api/v1/auth/check-default-password
func (h *AuthHandler) CheckDefaultPassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}

	isDefault, err := h.service.IsDefaultPassword(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"isDefaultPassword": isDefault,
		},
	})
}

// CreateUser 创建用户（管理员接口）
// POST /api/v1/auth/users
func (h *AuthHandler) CreateUser(c *gin.Context) {
	// 检查当前用户是否为管理员
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限，仅管理员可创建用户",
			"data":    nil,
		})
		return
	}

	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户创建成功",
		"data":    user.ToUserInfo(),
	})
}

// DeleteUser 删除用户（管理员接口）
// DELETE /api/v1/auth/users/:id
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	// 检查当前用户是否为管理员
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限，仅管理员可删除用户",
			"data":    nil,
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	if err := h.service.DeleteUser(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户删除成功",
		"data":    nil,
	})
}

// UpdateUser 更新用户信息（管理员接口）
// PUT /api/v1/auth/users/:id
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	// 检查当前用户是否为管理员
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限，仅管理员可更新用户",
			"data":    nil,
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	if err := h.service.UpdateUser(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户更新成功",
		"data":    nil,
	})
}

// GetUserList 获取用户列表（管理员接口，分页）
// GET /api/v1/auth/users?page=1&pageSize=20
func (h *AuthHandler) GetUserList(c *gin.Context) {
	// 检查当前用户是否为管理员
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限，仅管理员可查看用户列表",
			"data":    nil,
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	users, total, err := h.service.GetUserList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败",
			"data":    nil,
		})
		return
	}

	// 转换为用户信息（不包含密码）
	var userList []map[string]interface{}
	for _, user := range users {
		userList = append(userList, user.ToUserInfo())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    userList,
		"total":   total,
		"current": page,
	})
}
