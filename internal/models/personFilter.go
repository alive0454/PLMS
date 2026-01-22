package models

type PersonFilter struct {
	//楼号
	BuildingNumber string `json:"buildingNumber"`
	//单元
	UnitNumber int `json:"unitNumber"`
	//房号
	RoomNumber string `json:"roomNumber"`
	//姓名
	Name string `json:"name"`
	//身份证
	IDCard string `json:"idCard"`
	//年龄区间
	Age []string `json:"age"`
	//性别
	Gender int `json:"gender"`
	//是否常驻
	IsPermanent int `json:"isPermanent"`
	//住房情况
	HouseSituation string `json:"houseSituation"`
	//住房性质
	PropertyNature string `json:"propertyNature"`
	//户籍类型
	RegisteredResidenceType int `json:"registeredResidenceType"`
	//户籍地
	RegisteredResidence string `json:"registeredResidence"`
	//电话
	Telephone string `json:"telephone"`
	//是否有电动车
	HasElectricCar int `json:"hasElectricCar"`
	//失能等级
	DisabilityLevel string `json:"disabilityLevel"`
	//是否低保
	IsLowIncome int `json:"isLowIncome"`
	//是否低收入
	IsLowIncome2 int `json:"isLowIncome2"`
	//是否特困供养
	IsDestitute int `json:"isDestitute"`
	//是否计划生育特殊家庭
	IsFamilyPlanningSpecial int `json:"isFamilyPlanningSpecial"`
	//残疾类别及等级
	DisabilityCategory string `json:"disabilityCategory"`
	//是否独居
	IsLivingAlone int `json:"isLivingAlone"`
	//是否空巢
	IsEmptyNest int `json:"isEmptyNest"`
	//是否孤寡
	IsOrphaned int `json:"isOrphaned"`
	//是否需重点关注
	IsNeedsFocus int `json:"isNeedsFocus"`
	//电动车车牌号
	LicensePlate string `json:"licensePlate"`
	//是否入群
	IsInGroup int `json:"isInGroup"`
	//是否私信
	IsPrivateMessage int `json:"isPrivateMessage"`
	//是否有宠物
	HasPet int `json:"hasPet"`
	//查询类型交集（and）并集（or）
	QueryType int `json:"queryType"`
	Page      int `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize  int `json:"pageSize" form:"pageSize" binding:"omitempty,min=1,max=100000"`
}
