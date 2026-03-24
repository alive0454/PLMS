package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"PLMS/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWT密钥（从环境变量读取，默认为示例密钥）
var jwtSecret []byte

// 初始化JWT密钥
func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)
}

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService 认证服务
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务实例
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string                 `json:"token"`
	User  map[string]interface{} `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查询用户
	var user models.SysUser
	if err := s.db.Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginTime = &now
	s.db.Model(&user).Update("last_login_time", now)

	// 生成JWT token
	token, err := GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user.ToUserInfo(),
	}, nil
}

// GetUserByID 根据ID获取用户信息
func (s *AuthService) GetUserByID(userID int64) (*models.SysUser, error) {
	var user models.SysUser
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, req *ChangePasswordRequest) error {
	// 查询用户
	var user models.SysUser
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码和默认密码标志
	return s.db.Model(&user).Updates(map[string]interface{}{
		"password":            string(hashedPassword),
		"is_default_password": 0,
	}).Error
}

// IsDefaultPassword 检查是否使用默认密码
func (s *AuthService) IsDefaultPassword(userID int64) (bool, error) {
	var user models.SysUser
	if err := s.db.First(&user, userID).Error; err != nil {
		return false, err
	}
	return user.IsDefaultPassword == 1, nil
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64, username, role string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword 生成密码哈希（用于初始化用户）
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// CreateUser 创建新用户
func (s *AuthService) CreateUser(req *CreateUserRequest) (*models.SysUser, error) {
	// 检查用户名是否已存在
	var existingUser models.SysUser
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &models.SysUser{
		Username:          req.Username,
		Password:          hashedPassword,
		Name:              req.Name,
		Role:              req.Role,
		IsDefaultPassword: 0, // 新用户不是默认密码
		Status:            1, // 默认启用
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.New("创建用户失败: " + err.Error())
	}

	return user, nil
}

// DeleteUser 删除用户（软删除）
func (s *AuthService) DeleteUser(userID int64) error {
	result := s.db.Model(&models.SysUser{}).Where("id = ?", userID).Update("status", 0)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

// UpdateUser 更新用户信息
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
	Status   int8   `json:"status" binding:"omitempty,oneof=0 1"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

func (s *AuthService) UpdateUser(userID int64, req *UpdateUserRequest) error {
	updates := map[string]interface{}{
		"name": req.Name,
		"role": req.Role,
	}
	// 如果传了 status，也更新状态
	if req.Status == 0 || req.Status == 1 {
		updates["status"] = req.Status
	}
	// 如果传了密码，加密后更新
	if req.Password != "" {
		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
			return errors.New("密码加密失败")
		}
		updates["password"] = hashedPassword
		updates["is_default_password"] = 0 // 修改密码后不再是默认密码
	}

	result := s.db.Model(&models.SysUser{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

// GetUserList 获取用户列表（分页，包含所有状态用户，正常的排前面）
func (s *AuthService) GetUserList(page, pageSize int) ([]models.SysUser, int64, error) {
	var users []models.SysUser
	var total int64

	// 计算总数（包含所有状态）
	if err := s.db.Model(&models.SysUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按 status DESC 排序（1正常在前，0禁用在后）
	offset := (page - 1) * pageSize
	err := s.db.
		Order("status DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}
