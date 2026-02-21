package models

type PersonStatistic struct {
	TotalHouseholds            int64             `json:"total_households"`             // 总户数
	PermanentPopulation        int64             `json:"permanent_population"`         // 常住人口
	PermanentPopulationPercent string            `json:"permanent_population_percent"` // 常住人口占比
	FloatingPopulation         int64             `json:"floating_population"`          // 流动人口
	FloatingPopulationPercent  string            `json:"floating_population_percent"`  // 流动人口占比
	SelfOccupiedHouses         int64             `json:"self_occupied_houses"`         // 自住户数
	SelfOccupiedHousesPercent  string            `json:"self_occupied_houses_percent"` // 自住户占比
	RentedHouses               int64             `json:"rented_houses"`                // 出租户数
	RentedHousesPercent        string            `json:"rented_houses_percent"`        // 出租户占比
	VacantHouses               int64             `json:"vacant_houses"`                // 空置户数
	VacantHousesPercent        string            `json:"vacant_houses_percent"`        // 空置户占比
	DecorationHouses           int64             `json:"decoration_houses"`            // 装修户数
	DecorationHousesPercent    string            `json:"decoration_houses_percent"`    // 装修户占比
	RegisteredDist             map[string]string `json:"registered_dist"`              // 户籍分布
	AgeDist                    map[string]string `json:"age_dist"`                     // 年龄分布
	//HouseStatusDist           map[string]float64 `json:"house_status_dist"`            // 房屋状态分布
}
