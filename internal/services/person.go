package services

import (
	"PLMS/internal/models"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type PersonService struct {
	db *gorm.DB
}

func NewPersonService(db *gorm.DB) *PersonService {
	return &PersonService{db: db}
}

// GetBuildingNumbers 是 PersonService 的一个方法，用于获取所有不重复的楼号
// 返回一个字符串切片和可能的错误
func (p *PersonService) GetBuildingNumbers() ([]string, error) {
	// 初始化一个空的字符串切片用于存储建筑编号
	var buildingNumbers []string
	// 使用数据库查询，从 Person 表中获取不重复的 building_number 字段
	// 并将查询结果填充到 buildingNumbers 切片中
	err := p.db.Model(&models.Person{}).Where("is_del", 0).Distinct("building_number").Pluck("building_number", &buildingNumbers).Error
	// 返回建筑编号切片和可能的错误
	return buildingNumbers, err
}

// GetUnitNumbersByBuildingNumber 根据楼栋号获取所有单元号
// 参数:
//   - buildingNumber: 楼栋号字符串
//
// 返回值:
//   - []string: 包含所有单元号的字符串切片
//   - error: 可能发生的错误信息
func (p *PersonService) GetUnitNumbersByBuildingNumber(buildingNumber string) ([]int, error) {
	var unitNumbers []int // 用于存储查询结果的单元号切片
	// 执行数据库查询，从Person表中查询指定楼栋号的所有不重复的单元号
	err := p.db.Model(&models.Person{}).Where("is_del", 0).Where("building_number = ?", buildingNumber).Distinct("unit_number").Pluck("unit_number", &unitNumbers).Error
	return unitNumbers, err // 返回查询结果和可能的错误
}

// PersonService 结构体的方法，用于构建人员查询条件
// 参数:
//   - query: 初始的 GORM 查询对象
//   - filter: 人员筛选条件模型
//
// 返回值:
//   - *gorm.DB: 构建完成后的 GORM 查询对象
func (p *PersonService) buildPersonQuery(query *gorm.DB, filter models.PersonFilter) *gorm.DB {
	// 先处理 is_del 条件
	query = query.Where("is_del = ?", 0)

	// 如果 QueryType == 1，使用 OR 组合所有条件
	if filter.QueryType == 1 {
		// 收集所有 OR 条件
		var orConditions []clause.Expression

		// 楼号条件查询
		if filter.BuildingNumber != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "building_number"},
				Value:  "%" + filter.BuildingNumber + "%",
			})
		}

		// 单元号条件查询
		if filter.UnitNumber != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "unit_number"},
				Value:  filter.UnitNumber,
			})
		}

		// 房号条件查询
		if filter.RoomNumber != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "room_number"},
				Value:  "%" + filter.RoomNumber + "%",
			})
		}

		// 姓名条件查询
		if filter.Name != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "name"},
				Value:  "%" + filter.Name + "%",
			})
		}

		// 身份证号条件查询
		if filter.IDCard != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "id_card"},
				Value:  "%" + filter.IDCard + "%",
			})
		}

		// 年龄范围条件查询
		if len(filter.Age) > 2 {
			orConditions = append(orConditions, clause.And(
				clause.Gte{Column: clause.Column{Name: "age"}, Value: filter.Age[0]},
				clause.Lte{Column: clause.Column{Name: "age"}, Value: filter.Age[1]},
			))
		}

		// 性别条件查询
		if filter.Gender != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "gender"},
				Value:  filter.Gender,
			})
		}

		// 是否常住条件查询
		if filter.IsPermanent != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_permanent"},
				Value:  filter.IsPermanent,
			})
		}

		// 住房情况条件查询
		if filter.HouseSituation != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "house_situation"},
				Value:  "%" + filter.HouseSituation + "%",
			})
		}

		// 房产性质条件查询
		if filter.PropertyNature != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "property_nature"},
				Value:  "%" + filter.PropertyNature + "%",
			})
		}

		// 户籍类型条件查询
		if filter.RegisteredResidenceType != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "registered_residence_type"},
				Value:  filter.RegisteredResidenceType,
			})
		}

		if filter.RegisteredResidence != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "registered_residence"},
				Value:  "%" + filter.RegisteredResidence + "%",
			})
		}

		if filter.Telephone != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "telephone"},
				Value:  "%" + filter.Telephone + "%",
			})
		}

		if filter.HasElectricCar != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "has_electric_car"},
				Value:  filter.HasElectricCar,
			})
		}

		if filter.DisabilityLevel != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "disability_level"},
				Value:  "%" + filter.DisabilityLevel + "%",
			})
		}

		if filter.IsLowIncome != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_low_income"},
				Value:  filter.IsLowIncome,
			})
		}

		if filter.IsLowIncome2 != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_low_income2"},
				Value:  filter.IsLowIncome2,
			})
		}

		if filter.IsDestitute != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_destitute"},
				Value:  filter.IsDestitute,
			})
		}

		if filter.IsFamilyPlanningSpecial != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_family_planning_special"},
				Value:  filter.IsFamilyPlanningSpecial,
			})
		}

		if filter.DisabilityCategory != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "disability_category"},
				Value:  "%" + filter.DisabilityCategory + "%",
			})
		}

		if filter.IsLivingAlone != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_living_alone"},
				Value:  filter.IsLivingAlone,
			})
		}

		if filter.IsEmptyNest != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_empty_nest"},
				Value:  filter.IsEmptyNest,
			})
		}

		if filter.IsOrphaned != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_orphaned"},
				Value:  filter.IsOrphaned,
			})
		}

		if filter.IsNeedsFocus != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_needs_focus"},
				Value:  filter.IsNeedsFocus,
			})
		}

		if filter.LicensePlate != "" {
			orConditions = append(orConditions, clause.Like{
				Column: clause.Column{Name: "license_plate"},
				Value:  "%" + filter.LicensePlate + "%",
			})
		}

		if filter.IsInGroup != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_in_group"},
				Value:  filter.IsInGroup,
			})
		}

		if filter.IsPrivateMessage != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "is_private_message"},
				Value:  filter.IsPrivateMessage,
			})
		}

		if filter.HasPet != 0 {
			orConditions = append(orConditions, clause.Eq{
				Column: clause.Column{Name: "has_pet"},
				Value:  filter.HasPet,
			})
		}

		// 如果有 OR 条件，用 OR 连接
		if len(orConditions) > 0 {
			// 使用 clause.Or 连接所有条件
			query = query.Clauses(clause.Or(orConditions...))
		}
	} else {
		// QueryType != 0，使用 AND
		// 楼号条件查询
		if filter.BuildingNumber != "" {
			query = query.Where("building_number LIKE ?", "%"+filter.BuildingNumber+"%")
		}

		// 单元号条件查询
		if filter.UnitNumber != 0 {
			query = query.Where("unit_number = ?", filter.UnitNumber)
		}

		// 房号条件查询
		if filter.RoomNumber != "" {
			query = query.Where("room_number LIKE ?", "%"+filter.RoomNumber+"%")
		}

		// 姓名条件查询
		if filter.Name != "" {
			query = query.Where("name LIKE ?", "%"+filter.Name+"%")
		}

		// 身份证号条件查询
		if filter.IDCard != "" {
			query = query.Where("id_card LIKE ?", "%"+filter.IDCard+"%")
		}

		// 年龄范围条件查询
		if len(filter.Age) > 2 {
			query = query.Where("age >= ? AND age <= ?", filter.Age[0], filter.Age[1])
		}

		// 性别条件查询
		if filter.Gender != 0 {
			query = query.Where("gender = ?", filter.Gender)
		}

		// 是否常住条件查询
		if filter.IsPermanent != 0 {
			query = query.Where("is_permanent = ?", filter.IsPermanent)
		}

		// 住房情况条件查询
		if filter.HouseSituation != "" {
			query = query.Where("house_situation LIKE ?", "%"+filter.HouseSituation+"%")
		}

		// 房产性质条件查询
		if filter.PropertyNature != "" {
			query = query.Where("property_nature LIKE ?", "%"+filter.PropertyNature+"%")
		}

		// 户籍类型条件查询
		if filter.RegisteredResidenceType != 0 {
			query = query.Where("registered_residence_type = ?", filter.RegisteredResidenceType)
		}

		if filter.RegisteredResidence != "" {
			query = query.Where("registered_residence LIKE ?", "%"+filter.RegisteredResidence+"%")
		}

		if filter.Telephone != "" {
			query = query.Where("telephone LIKE ?", "%"+filter.Telephone+"%")
		}

		if filter.HasElectricCar != 0 {
			query = query.Where("has_electric_car = ?", filter.HasElectricCar)
		}

		if filter.DisabilityLevel != "" {
			query = query.Where("disability_level LIKE ?", "%"+filter.DisabilityLevel+"%")
		}

		if filter.IsLowIncome != 0 {
			query = query.Where("is_low_income = ?", filter.IsLowIncome)
		}

		if filter.IsLowIncome2 != 0 {
			query = query.Where("is_low_income2 = ?", filter.IsLowIncome2)
		}

		if filter.IsDestitute != 0 {
			query = query.Where("is_destitute = ?", filter.IsDestitute)
		}

		if filter.IsFamilyPlanningSpecial != 0 {
			query = query.Where("is_family_planning_special = ?", filter.IsFamilyPlanningSpecial)
		}

		if filter.DisabilityCategory != "" {
			query = query.Where("disability_category LIKE ?", "%"+filter.DisabilityCategory+"%")
		}

		if filter.IsLivingAlone != 0 {
			query = query.Where("is_living_alone = ?", filter.IsLivingAlone)
		}

		if filter.IsEmptyNest != 0 {
			query = query.Where("is_empty_nest = ?", filter.IsEmptyNest)
		}

		if filter.IsOrphaned != 0 {
			query = query.Where("is_orphaned = ?", filter.IsOrphaned)
		}

		if filter.IsNeedsFocus != 0 {
			query = query.Where("is_needs_focus = ?", filter.IsNeedsFocus)
		}

		if filter.LicensePlate != "" {
			query = query.Where("license_plate LIKE ?", "%"+filter.LicensePlate+"%")
		}

		if filter.IsInGroup != 0 {
			query = query.Where("is_in_group = ?", filter.IsInGroup)
		}

		if filter.IsPrivateMessage != 0 {
			query = query.Where("is_private_message = ?", filter.IsPrivateMessage)
		}

		if filter.HasPet != 0 {
			query = query.Where("has_pet = ?", filter.HasPet)
		}
	}

	return query
}
func (p *PersonService) buildPageQuery(query *gorm.DB, filter models.PersonFilter) *gorm.DB {
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	} else {
		// 如果没有分页参数，设置默认值
		filter.Page = 1
		filter.PageSize = 20
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}
	return query
}

func (p *PersonService) GetPersonInfo(id int) (models.PersonInfo, error) {
	var person models.Person
	var bicycles []models.ElectricBicycle
	var personInfo models.PersonInfo
	result := p.db.First(&person, id)
	personInfo.Person = person
	result = p.db.Model(&models.ElectricBicycle{}).Where("person_id=?", id).Find(&bicycles)
	personInfo.Bicycles = bicycles
	return personInfo, result.Error
}

func (p *PersonService) GetPersonInfoByRoom(buildingNumber, unitNumber, roomNumber string) ([]models.PersonInfo, error) {
	var persons []models.Person
	var bicycles []models.ElectricBicycle
	var personInfos []models.PersonInfo
	result := p.db.Where("building_number=? and unit_number=? and room_number=?",
		buildingNumber, unitNumber, roomNumber).Find(&persons)
	for _, person := range persons {
		personInfo := models.PersonInfo{}
		personInfo.Person = person
		result = p.db.Model(&models.ElectricBicycle{}).Where("person_id=?", person.ID).Find(&bicycles)
		personInfo.Bicycles = bicycles
		personInfos = append(personInfos, personInfo)
	}
	return personInfos, result.Error
}

// GetPersons 获取人员列表，支持过滤和分页
// 参数:
//   - filter: 人员过滤条件，包含查询条件和分页信息
//
// 返回值:
//   - []models.Person: 人员列表
//   - int64: 总记录数
//   - error: 错误信息
func (p *PersonService) GetPersons(filter models.PersonFilter) ([]models.Person, int64, error) {
	var persons []models.Person                             // 存储查询结果的人员列表
	var total int64                                         // 存储总记录数
	query := p.db.Model(&models.Person{})                   // 创建基础查询对象
	query = p.buildPersonQuery(query, filter)               // 构建查询条件
	query.Order("building_number,unit_number, room_number") // 设置排序规则，按楼号、单元号、房间号排序
	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err // 如果计数出错，返回错误
	}
	// 添加分页
	query = p.buildPageQuery(query, filter) // 根据过滤条件添加分页

	result := query.Find(&persons)      // 执行查询
	return persons, total, result.Error // 返回查询结果、总记录数和可能的错误
}

func (p *PersonService) GetRooms(filter models.PersonFilter) ([]models.RoomList, int64, error) {
	var roomList = []models.RoomList{}
	var total int64
	query := p.db.Table("person").
		Select("building_number, unit_number, room_number")
	query = p.buildPersonQuery(query, filter)
	query.Where("building_number<>'' and unit_number>0 and room_number<>''")
	query.Group("building_number, unit_number, room_number")
	query.Order("building_number, unit_number, room_number")
	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err // 如果计数出错，返回错误
	}
	query = p.buildPageQuery(query, filter) // 根据过滤条件添加分页
	result := query.Find(&roomList)         // 执行查询
	var persons = []models.Person{}
	query = p.db.Model(&models.Person{})
	query.Where("is_del=0")

	var orConditions []clause.Expression
	for _, room := range roomList {
		orConditions = append(orConditions, clause.And(
			clause.Eq{Column: "building_number", Value: room.BuildingNumber},
			clause.Eq{Column: "unit_number", Value: room.UnitNumber},
			clause.Eq{Column: "room_number", Value: room.RoomNumber},
		))
	}
	query.Where(clause.Or(orConditions...))
	query.Find(&persons)
	for i := range roomList {
		room := &roomList[i]
		//人员信息 如果人员里面有租户，就写第一个租户的名字，如果没有，就看有没有包含业的第一个名字，如果都没有，就歇第一个名字
		var rentData = []models.Person{}
		linq.From(persons).Where(func(p interface{}) bool {
			person := p.(models.Person)
			return person.BuildingNumber == room.BuildingNumber && person.UnitNumber == room.UnitNumber && person.RoomNumber == room.RoomNumber
		}).Where(func(p interface{}) bool {
			return strings.Contains(p.(models.Person).HousingSituation, "租")
		}).OrderBy(func(p interface{}) interface{} {
			return p.(models.Person).ID
		}).ToSlice(&rentData)
		var ownerData = []models.Person{}
		linq.From(persons).Where(func(p interface{}) bool {
			person := p.(models.Person)
			return person.BuildingNumber == room.BuildingNumber && person.UnitNumber == room.UnitNumber && person.RoomNumber == room.RoomNumber
		}).Where(func(p interface{}) bool {
			return strings.Contains(p.(models.Person).HousingSituation, "业")
		}).OrderBy(func(p interface{}) interface{} {
			return p.(models.Person).ID
		}).ToSlice(&ownerData)
		var emptyData = []models.Person{}
		linq.From(persons).Where(func(p interface{}) bool {
			person := p.(models.Person)
			return person.BuildingNumber == room.BuildingNumber && person.UnitNumber == room.UnitNumber && person.RoomNumber == room.RoomNumber
		}).Where(func(p interface{}) bool {
			return strings.Contains(p.(models.Person).HousingSituation, "空")
		}).OrderBy(func(p interface{}) interface{} {
			return p.(models.Person).ID
		}).ToSlice(&emptyData)
		if len(rentData) > 0 {
			room.ContactName = rentData[0].Name
			room.PersonNum = len(rentData)
			room.Telephone = rentData[0].Telephone
			room.HousingSituation = "出租"
		} else if len(ownerData) > 0 {
			room.ContactName = ownerData[0].Name
			room.Telephone = ownerData[0].Telephone
			if len(emptyData) > 0 {
				room.PersonNum = 0
				room.HousingSituation = "空置"
			} else {
				room.PersonNum = len(ownerData)
				room.HousingSituation = "自住"
			}
		}
	}
	return roomList, total, result.Error
}

// PersonService 结构体的方法，获取人口统计信息
// 参数:
//   - buildingNumber: 楼号，为空或"0"时不限制楼号
//
// 返回值:
//   - *models.PersonStatistic: 人口统计结果
//   - error: 错误信息
func (p *PersonService) GetPersonStatistics(buildingNumber string) (*models.PersonStatistic, error) {
	var result models.PersonStatistic
	// SQL查询语句，用于统计各类人口信息
	sql := `
    SELECT 
        COUNT(DISTINCT room_number) AS total_households,                    -- 总户数
        SUM(CASE WHEN is_permanent = 1 THEN 1 ELSE 0 END) AS permanent_population,    -- 常住人口
        SUM(CASE WHEN is_permanent = 2 THEN 1 ELSE 0 END) AS floating_population,      -- 流动人口
        COUNT(DISTINCT CASE WHEN property_nature LIKE ? THEN room_number END) AS self_occupied_houses,  -- 自住房数量
        COUNT(DISTINCT CASE WHEN housing_situation LIKE ? THEN room_number END) AS rented_houses,       -- 租住房数量
        COUNT(DISTINCT CASE WHEN housing_situation LIKE ? THEN room_number END) AS vacant_houses,        -- 空置房数量
        COUNT(DISTINCT CASE WHEN housing_situation LIKE ? THEN room_number END) AS decoration_houses    -- 装修房数量
    FROM person
    WHERE is_del = 0  -- 只查询未删除的记录
    `

	// SQL查询参数
	params := []interface{}{
		"%自住%", // property_nature LIKE ? 的参数
		"%租%",  // housing_situation LIKE ? 的参数
		"%空%",  // housing_situation LIKE ? 的参数
		"%装修%", // housing_situation LIKE ? 的参数
	}
	// 动态添加楼号条件
	if buildingNumber != "" && buildingNumber != "0" {
		sql += " AND building_number = ?" // 添加楼号查询条件
		//params = append([]interface{}{buildingNumber}, params...)
		params = append(params, buildingNumber) // 添加楼号参数
		fmt.Println(params)
	}
	// 执行SQL查询并将结果扫描到result结构体中
	err := p.db.Raw(sql, params...).Scan(&result).Error
	if err != nil {
		return nil, err // 查询出错时返回错误
	}
	//统计信息
	totalPersonCount := result.PermanentPopulation + result.FloatingPopulation
	if result.TotalHouseholds > 0 {
		//流动人口比例
		result.FloatingPopulationPercent = fmt.Sprintf("%.2f%%", float64(result.FloatingPopulation)/
			float64(totalPersonCount)*100)
		//自住房比例
		result.SelfOccupiedHousesPercent = fmt.Sprintf("%.2f%%", float64(result.SelfOccupiedHouses)/float64(result.TotalHouseholds)*100)
		//租住房比例
		result.RentedHousesPercent = fmt.Sprintf("%.2f%%", float64(result.RentedHouses)/float64(result.TotalHouseholds)*100)
		//空置房比例
		result.VacantHousesPercent = fmt.Sprintf("%.2f%%", float64(result.VacantHouses)/float64(result.TotalHouseholds)*100)
		//装修房比例
		result.DecorationHousesPercent = fmt.Sprintf("%.2f%%", float64(result.DecorationHouses)/float64(result.TotalHouseholds)*100)
	}

	stat2, err := p.getPersonDemographicStats(buildingNumber)
	if err == nil {
		result.RegisteredDist = stat2.RegisteredResidenceTypeStats
		result.AgeDist = stat2.AgeDist
		fmt.Println(result)
	}

	return &result, nil // 返回查询结果
}
func (p *PersonService) getPersonDemographicStats(buildingNumber string) (*models.PersonDemographicStat, error) {
	var result models.PersonDemographicStat

	// 初始化map
	result.RegisteredResidenceTypeStats = make(map[string]string)
	result.AgeDist = make(map[string]string)

	// 构建SQL查询
	sql := `
    SELECT 
        -- 户籍类型统计
        SUM(CASE WHEN registered_residence_type = 1 THEN 1 ELSE 0 END) as type1,
        SUM(CASE WHEN registered_residence_type = 2 THEN 1 ELSE 0 END) as type2,
        SUM(CASE WHEN registered_residence_type = 3 THEN 1 ELSE 0 END) as type3,
        SUM(CASE WHEN registered_residence_type NOT IN (1,2,3) THEN 1 ELSE 0 END) as type_other,
        
        -- 独居统计
        SUM(CASE WHEN is_living_alone = 1 THEN 1 ELSE 0 END) as living_alone,
        
        -- 空巢统计
        SUM(CASE WHEN is_empty_nest = 1 THEN 1 ELSE 0 END) as empty_nest,
        
        -- 年龄段统计
        SUM(CASE WHEN age >= 60 AND age <= 80 THEN 1 ELSE 0 END) as age60_to80,
        SUM(CASE WHEN age > 80 THEN 1 ELSE 0 END) as age_over80,
        
        -- 总人数
        COUNT(*) as total
    FROM person
    WHERE is_del = 0
    `

	var params []interface{}
	if buildingNumber != "" && buildingNumber != "0" {
		sql += " AND building_number = ?"
		params = append(params, buildingNumber)
	}

	var stats struct {
		Type1       int64
		Type2       int64
		Type3       int64
		TypeOther   int64
		LivingAlone int64
		EmptyNest   int64
		Age60To80   int64
		AgeOver80   int64
		Total       int64
	}

	err := p.db.Raw(sql, params...).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	result.TotalPopulation = stats.Total
	// 填充结果
	result.RegisteredResidenceTypeStats["东湖户籍"] = fmt.Sprintf("%.2f%%", float64(stats.Type1)/float64(result.TotalPopulation)*100)
	result.RegisteredResidenceTypeStats["北京市其他户籍"] = fmt.Sprintf("%.2f%%", float64(stats.Type2)/float64(result.TotalPopulation)*100)
	result.RegisteredResidenceTypeStats["外地户籍"] = fmt.Sprintf("%.2f%%", float64(stats.Type3)/float64(result.TotalPopulation)*100)
	result.RegisteredResidenceTypeStats["其他户籍"] = fmt.Sprintf("%.2f%%", float64(stats.TypeOther)/float64(result.TotalPopulation)*100)

	result.AgeDist["空巢"] = fmt.Sprintf("%.2f%%", float64(stats.EmptyNest)/float64(result.TotalPopulation)*100)
	result.AgeDist["独居"] = fmt.Sprintf("%.2f%%", float64(stats.LivingAlone)/float64(result.TotalPopulation)*100)
	result.AgeDist["60岁及以上"] = fmt.Sprintf("%.2f%%", float64(stats.Age60To80)/float64(result.TotalPopulation)*100)
	result.AgeDist["80岁及以上"] = fmt.Sprintf("%.2f%%", float64(stats.AgeOver80)/float64(result.TotalPopulation)*100)

	return &result, nil
}
