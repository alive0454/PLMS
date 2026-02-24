package models

import (
	"time"
)

// SysUser 系统用户表
type SysUser struct {
	ID                int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	Username          string     `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Password          string     `json:"-" gorm:"size:255;not null;comment:密码（BCrypt加密）"`
	Name              string     `json:"name" gorm:"size:50;comment:真实姓名"`
	Role              string     `json:"role" gorm:"size:50;default:'user';comment:角色：admin/user"`
	IsDefaultPassword int8       `json:"isDefaultPassword" gorm:"column:is_default_password;type:tinyint(1);default:1;comment:是否使用默认密码：1是，0否"`
	Status            int8       `json:"status" gorm:"type:tinyint(1);default:1;comment:账号状态：1启用，0禁用"`
	LastLoginTime     *time.Time `json:"lastLoginTime" gorm:"column:last_login_time;comment:最后登录时间"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}

// ToUserInfo 转换为用户信息（用于返回给前端，不包含敏感信息）
func (u *SysUser) ToUserInfo() map[string]interface{} {
	avatar := ""
	if u.Role == "admin" {
		avatar = "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png"
	} else {
		avatar = "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png"
	}
	return map[string]interface{}{
		"id":                u.ID,
		"username":          u.Username,
		"name":              u.Name,
		"role":              u.Role,
		"isDefaultPassword": u.IsDefaultPassword == 1,
		"avatar":            avatar,
	}
}
