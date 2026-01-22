package models

import (
	"time"
)

// ElectricBicycle 电动车信息表
type ElectricBicycle struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                      // 主键ID
	PersonID    int64     `gorm:"column:person_id;not null" json:"person_id"`                        // 所属人员ID
	Model       string    `gorm:"column:model;type:varchar(50);not null" json:"model"`               // 车型型号
	Brand       *string   `gorm:"column:brand;type:varchar(30)" json:"brand"`                        // 品牌
	Color       *string   `gorm:"column:color;type:varchar(20)" json:"color"`                        // 颜色
	PlateNumber string    `gorm:"column:plate_number;type:varchar(20);not null" json:"plate_number"` // 车牌号码
	IsDel       int8      `gorm:"column:is_del;default:0" json:"is_del"`                             // 软删除标记
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`                // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                // 更新时间
}

// TableName 指定表名
func (ElectricBicycle) TableName() string {
	return "electric_bicycle"
}
