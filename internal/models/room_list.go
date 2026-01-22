package models

type RoomList struct {
	//楼号
	BuildingNumber string `json:"buildingNumber"`
	//单元
	UnitNumber int `json:"unitNumber"`
	//房号
	RoomNumber string `json:"roomNumber"`
	//联系人
	ContactName string `json:"contactName"`
	//居住人数
	PersonNum int `json:"personNum"`
	//联系电话
	Telephone string `json:"telephone"`
	//住房情况
	HousingSituation string `json:"housingSituation"`
	//总数
	Total int `json:"total"`
}
