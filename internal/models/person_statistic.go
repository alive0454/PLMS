package models

// 在 models/person_statistic.go 中添加
type PersonDemographicStat struct {
	// 户籍类型统计
	RegisteredResidenceTypeStats map[string]string `json:"registered_residence_type_stats"` // 户籍类型统计
	AgeDist                      map[string]string `json:"age_dist"`                        // 年龄分布
	//HouseStatusDist              map[string]float64 `json:"house_status_dist"`               // 房屋状态分布

	// 可以添加总计用于计算百分比
	TotalPopulation int64 `json:"total_population"` // 总人数
}
