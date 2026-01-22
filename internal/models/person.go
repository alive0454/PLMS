package models

import (
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
	IsDel                   int       `gorm:"column:is_del;type:tinyint;default:0" json:"is_del"`
	CreatedAt               time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (Person) TableName() string {
	return "person"
}
