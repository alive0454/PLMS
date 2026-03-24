package models

import (
	"strconv"
	"time"
)

// Person 人员信息台账
type Person struct {
	ID                      int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BuildingNumber          string    `gorm:"column:building_number;not null;type:varchar(20)" json:"building_number"`
	UnitNumber              int       `gorm:"column:unit_number;type:int" json:"unit_number"`
	RoomNumber              string    `gorm:"column:room_number;not null;type:varchar(100)" json:"room_number"`
	Name                    string    `gorm:"column:name;type:varchar(50)" json:"name"`
	IDCard                  string    `gorm:"column:id_card;type:varchar(20);index" json:"id_card"`
	Age                     int       `gorm:"column:age;type:int;index" json:"age"`
	Gender                  int       `gorm:"column:gender;type:tinyint" json:"gender"`
	IsPermanent             int       `gorm:"column:is_permanent;type:tinyint" json:"is_permanent"`
	HousingSituation        string    `gorm:"column:housing_situation;type:varchar(100)" json:"housing_situation"`
	PropertyNature          string    `gorm:"column:property_nature;type:varchar(200)" json:"property_nature"`
	RegisteredResidenceType int       `gorm:"column:registered_residence_type;type:tinyint" json:"registered_residence_type"`
	RegisteredResidence     string    `gorm:"column:registered_residence;type:varchar(200)" json:"registered_residence"`
	Telephone               string    `gorm:"column:telephone;type:varchar(50)" json:"telephone"`
	FirstContact            string    `gorm:"column:first_contact;type:varchar(50)" json:"first_contact"`
	ElderRelationship       string    `gorm:"column:elder_relationship;type:varchar(50)" json:"elder_relationship"`
	ElderContactPhone       string    `gorm:"column:elder_contact_phone;type:varchar(20)" json:"elder_contact_phone"`
	SpecialSituation        string    `gorm:"column:special_situation;type:text" json:"special_situation"`
	HasElectricCar          int       `gorm:"column:has_electric_car;type:tinyint" json:"has_electric_car"`
	DisabilityLevel         string    `gorm:"column:disability_level;type:varchar(100)" json:"disability_level"`
	IsLowIncome             int       `gorm:"column:is_low_income;type:tinyint" json:"is_low_income"`
	IsLowIncome2            int       `gorm:"column:is_low_income2;type:tinyint" json:"is_low_income2"`
	IsDestitute             int       `gorm:"column:is_destitute;type:tinyint" json:"is_destitute"`
	IsFamilyPlanningSpecial int       `gorm:"column:is_family_planning_special;type:tinyint" json:"is_family_planning_special"`
	DisabilityCategory      string    `gorm:"column:disability_category;type:varchar(100)" json:"disability_category"`
	IsLivingAlone           int       `gorm:"column:is_living_alone;type:tinyint" json:"is_living_alone"`
	IsEmptyNest             int       `gorm:"column:is_empty_nest;type:tinyint" json:"is_empty_nest"`
	IsOrphaned              int       `gorm:"column:is_orphaned;type:tinyint" json:"is_orphaned"`
	IsNeedsFocus            int       `gorm:"column:is_needs_focus;type:tinyint" json:"is_needs_focus"`
	OtherSituation          string    `gorm:"column:other_situation;type:text" json:"other_situation"`
	LicensePlate            string    `gorm:"column:license_plate;type:varchar(20)" json:"license_plate"`
	BrandModel              string    `gorm:"column:brand_model;type:varchar(100)" json:"brand_model"`
	IsInGroup               int       `gorm:"column:is_in_group;type:tinyint" json:"is_in_group"`
	IsPrivateMessage        int       `gorm:"column:is_private_message;type:tinyint" json:"is_private_message"`
	HasPet                  int       `gorm:"column:has_pet;type:tinyint" json:"has_pet"`
	LastContactTime         string    `gorm:"column:last_contact_time;type:varchar(500)" json:"last_contact_time"`
	OtherInfo               string    `gorm:"column:other_info;type:text" json:"other_info"`
	IsCp                    int       `gorm:"column:is_cp;type:tinyint;default:0" json:"is_cp"`       // 是否党员：0否，1是
	CpJoiningDay            *string   `gorm:"column:cp_joining_day;type:date" json:"cp_joining_day"`  // 入党日
	Nationality             string    `gorm:"column:nationality;type:varchar(50)" json:"nationality"` // 民族
	Education               string    `gorm:"column:education;type:varchar(50)" json:"education"`     // 学历
	CpRemark                string    `gorm:"column:cp_remark;type:text" json:"cp_remark"`            // 党员备注
	IsDel                   int       `gorm:"column:is_del;type:tinyint;default:0" json:"is_del"`
	CreatedAt               time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (Person) TableName() string {
	return "person"
}

// GetExportValue 获取字段导出值（带转换）
func (p *Person) GetExportValue(field string) string {
	switch field {
	case "building_number":
		return p.BuildingNumber
	case "unit_number":
		return strconv.Itoa(p.UnitNumber)
	case "room_number":
		return p.RoomNumber
	case "name":
		return p.Name
	case "id_card":
		return p.IDCard
	case "age":
		return strconv.Itoa(p.Age)
	case "gender":
		switch p.Gender {
		case 1:
			return "男"
		case 2:
			return "女"
		default:
			return "未知"
		}
	case "telephone":
		return p.Telephone
	case "is_permanent":
		return convertYesNo(p.IsPermanent)
	case "housing_situation":
		return p.HousingSituation
	case "property_nature":
		return p.PropertyNature
	case "registered_residence":
		return p.RegisteredResidence
	case "registered_residence_type":
		return convertResidenceType(p.RegisteredResidenceType)
	case "is_living_alone":
		return convertYesNo(p.IsLivingAlone)
	case "is_empty_nest":
		return convertYesNo(p.IsEmptyNest)
	case "has_electric_car":
		return convertYesNo(p.HasElectricCar)
	case "license_plate":
		return p.LicensePlate
	case "is_low_income":
		return convertYesNo(p.IsLowIncome)
	case "is_low_income2":
		return convertYesNo(p.IsLowIncome2)
	case "is_destitute":
		return convertYesNo(p.IsDestitute)
	case "has_pet":
		return convertYesNo(p.HasPet)
	case "is_cp":
		return convertYesNo(p.IsCp)
	case "nationality":
		return p.Nationality
	case "education":
		return p.Education
	case "cp_joining_day":
		if p.CpJoiningDay != nil {
			return *p.CpJoiningDay
		}
		return ""
	case "cp_remark":
		return p.CpRemark
	case "first_contact":
		return p.FirstContact
	case "elder_relationship":
		return p.ElderRelationship
	case "elder_contact_phone":
		return p.ElderContactPhone
	case "special_situation":
		return p.SpecialSituation
	case "disability_level":
		return p.DisabilityLevel
	case "is_family_planning_special":
		return convertYesNo(p.IsFamilyPlanningSpecial)
	case "disability_category":
		return p.DisabilityCategory
	case "is_orphaned":
		return convertYesNo(p.IsOrphaned)
	case "is_needs_focus":
		return convertYesNo(p.IsNeedsFocus)
	case "other_situation":
		return p.OtherSituation
	case "brand_model":
		return p.BrandModel
	case "is_in_group":
		return convertYesNo(p.IsInGroup)
	case "is_private_message":
		return convertYesNo(p.IsPrivateMessage)
	case "last_contact_time":
		return p.LastContactTime
	case "other_info":
		return p.OtherInfo
	default:
		return ""
	}
}

// convertYesNo 转换是否类字段
func convertYesNo(val int) string {
	switch val {
	case 1:
		return "是"
	case 2:
		return "否"
	default:
		return "未知"
	}
}

// convertResidenceType 转换户籍类型
func convertResidenceType(val int) string {
	switch val {
	case 1:
		return "东湖户籍"
	case 2:
		return "北京市其他户籍"
	case 3:
		return "外地户籍"
	default:
		return "其他户籍"
	}
}

// GetExportFieldHeader 获取导出字段中文名
func GetExportFieldHeader(field string) string {
	headers := map[string]string{
		"building_number":            "楼号",
		"unit_number":                "单元",
		"room_number":                "房间号",
		"name":                       "姓名",
		"id_card":                    "身份证号",
		"age":                        "年龄",
		"gender":                     "性别",
		"telephone":                  "电话",
		"is_permanent":               "是否常驻",
		"housing_situation":          "住房情况",
		"property_nature":            "房屋性质",
		"registered_residence":       "户籍地",
		"registered_residence_type":  "户籍情况",
		"is_living_alone":            "是否独居",
		"is_empty_nest":              "是否空巢",
		"has_electric_car":           "是否有电动车",
		"license_plate":              "车牌号",
		"is_low_income":              "是否低保",
		"is_low_income2":             "是否低收入",
		"is_destitute":               "是否特困",
		"has_pet":                    "是否有宠物",
		"is_cp":                      "是否党员",
		"nationality":                "民族",
		"education":                  "学历",
		"cp_joining_day":             "入党日期",
		"cp_remark":                  "党员备注",
		"first_contact":              "第一联系人",
		"elder_relationship":         "与老人关系",
		"elder_contact_phone":        "紧急联系电话",
		"special_situation":          "特殊情况",
		"disability_level":           "失能等级",
		"is_family_planning_special": "是否计生特殊家庭",
		"disability_category":        "残疾类别及等级",
		"is_orphaned":                "是否孤寡",
		"is_needs_focus":             "是否重点关注",
		"other_situation":            "其他情况",
		"brand_model":                "电动车品牌型号",
		"is_in_group":                "是否入群",
		"is_private_message":         "是否私信",
		"last_contact_time":          "最后联系时间",
		"other_info":                 "其他信息",
	}
	if h, ok := headers[field]; ok {
		return h
	}
	return field
}
