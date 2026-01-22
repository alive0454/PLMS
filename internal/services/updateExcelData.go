package services

import (
	"PLMS/internal/models"
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type UpdateExcelDataService struct {
	db *gorm.DB
}

func NewUpdateExcelDataService(db *gorm.DB) *UpdateExcelDataService {
	return &UpdateExcelDataService{db: db}
}

func (s *UpdateExcelDataService) ImportExcelData(filePath string) error {
	processings := []string{"101楼", "103楼", "104楼新版", "105楼副本", "106楼副本", "108楼副本", "109楼-更新中", "110楼", "111楼新",
		"113楼", "114楼更新", "117楼", "118楼新版", "119楼", "120楼新版", "121楼新版", "122楼新版"}
	//processings := []string{"101楼"}
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取所有工作表
	sheets := f.GetSheetList()

	for _, sheet := range sheets {
		// 判断是否需要处理
		if !slices.Contains(processings, sheet) {
			continue
		}

		fmt.Printf("正在处理工作表: %s\n", sheet)
		// 读取人员数据（传入文件对象和工作表名称）
		persons, err := readPersonData(f, sheet)
		if err != nil {
			log.Printf("读取人员数据失败: %v", err)
			continue
		}
		fmt.Println("总数：" + strconv.Itoa(len(persons)))
		//保存人员数据
		//err = s.savePersons(persons)
		//if err != nil {
		//	log.Printf("保存人员信息失败: %v", err)
		//}
	}

	return nil
}

// getBuildingNumberFromMergedCells 专门处理B列的合并单元格
func getBuildingNumberFromMergedCells(f *excelize.File, sheet string, rows [][]string) (map[int]string, error) {
	// 获取所有合并单元格
	mergeCells, err := f.GetMergeCells(sheet)
	if err != nil {
		return nil, err
	}

	// 只处理B列（第二列）的合并单元格
	buildingMap := make(map[int]string)

	for _, mc := range mergeCells {
		startCell, endCell := mc.GetStartAxis(), mc.GetEndAxis()
		value := mc.GetCellValue()

		// 检查是否是B列的合并单元格
		startCol, startRow, _ := excelize.SplitCellName(startCell)
		endCol, endRow, _ := excelize.SplitCellName(endCell)

		// 只处理B列（列名是"B"）
		if startCol == "B" && endCol == "B" {
			// 将这个合并范围内的所有行都映射到同一个楼号值
			for row := startRow; row <= endRow; row++ {
				buildingMap[row] = value
			}
		}
	}

	return buildingMap, nil
}

/*
 * readPersonData 读取excel中人员信息
 * @param f: excelize.File类型的Excel文件对象
 * @param sheet: 要读取的工作表名称
 * @return: 返回Person结构体切片和可能的错误
 */
func readPersonData(f *excelize.File, sheet string) ([]models.Person, error) {
	var persons []models.Person

	// 获取所有行
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	// 获取B列合并单元格的映射
	buildingMap, err := getBuildingNumberFromMergedCells(f, sheet, rows)
	if err != nil {
		return nil, err
	}

	// 跳过表头行，从数据行开始读取
	startRow := 8
	for i := startRow; i < len(rows); i++ {
		var person models.Person

		// 解析楼号、单元、户号（处理B列的合并单元格）
		excelRowNum := i + 1 // Excel行号是从1开始的

		var buildingNumber string

		// 首先检查是否是合并单元格的一部分
		if mergedValue, exists := buildingMap[excelRowNum]; exists {
			// 如果是合并单元格，使用合并的值
			buildingNumber = mergedValue
		} else if len(rows[i]) > 1 && rows[i][1] != "" {
			// 如果不是合并单元格，直接读取B列的值
			buildingNumber = strings.TrimSpace(rows[i][1])
		}
		// 处理地址分割
		if buildingNumber != "" {
			//单独的房号
			if sheet == "108楼副本" || sheet == "110楼" || sheet == "111楼新" || sheet == "113楼" {
				bNum := extractNumbers(sheet)
				person.BuildingNumber = bNum
				person.UnitNumber = 1
				person.RoomNumber = buildingNumber
			} else if sheet == "117楼" { //117楼 楼、单元、房号
				building, unit, room := parseBuildingAddress(buildingNumber)
				person.BuildingNumber = building
				person.UnitNumber, _ = strconv.Atoi(unit)
				person.RoomNumber = room
			} else if sheet == "118楼新版" || sheet == "119楼" || sheet == "120楼新版" || sheet == "121楼新版" || sheet == "122楼新版" { //楼号-单元-房号
				splitAddress := strings.Split(buildingNumber, "-")
				person.BuildingNumber = splitAddress[0]
				person.UnitNumber, _ = strconv.Atoi(splitAddress[1])
				person.RoomNumber = splitAddress[2]
			} else { //楼号-房号
				parts := strings.SplitN(buildingNumber, "-", 2)
				if len(parts) == 2 {
					person.BuildingNumber = parts[0]
					person.UnitNumber = 1
					person.RoomNumber = parts[1]
				}
			}
		}
		// 解析姓名
		if len(rows[i]) > 2 {
			person.Name = sql.NullString{String: strings.TrimSpace(rows[i][2]), Valid: rows[i][2] != ""}.String
		}
		// 解析身份证
		if len(rows[i]) > 3 {
			person.IDCard = sql.NullString{String: strings.TrimSpace(rows[i][3]), Valid: rows[i][3] != ""}.String
		}

		// 解析年龄
		if len(rows[i]) > 4 && rows[i][4] != "" {
			num, err := strconv.Atoi(rows[i][4])
			if err == nil {
				person.Age = num
			}
		}

		// 解析性别
		if len(rows[i]) > 5 {
			gender := strings.TrimSpace(rows[i][5])
			if gender == "男" {
				person.Gender = 1
			} else if gender == "女" {
				person.Gender = 2
			}
		}

		//是否常驻
		//is_permanent tinyint COMMENT '是否常住：0未知，1是，2否',
		if len(rows[i]) > 6 {
			isPermanent := strings.TrimSpace(rows[i][6])
			if isPermanent == "是" {
				person.IsPermanent = 1 // 注意：这里原来是 person.Gender = 1，应该是笔误
			} else if isPermanent == "否" {
				person.IsPermanent = 2 // 注意：这里原来是 person.Gender = 2，应该是笔误
			}
		}

		// housing_situation VARCHAR(100) COMMENT '住房情况',
		if len(rows[i]) > 7 {
			person.HousingSituation = sql.NullString{String: strings.TrimSpace(rows[i][7]), Valid: rows[i][7] != ""}.String
		}

		//property_nature VARCHAR(200) COMMENT '房屋性质',
		if len(rows[i]) > 8 {
			person.PropertyNature = sql.NullString{String: strings.TrimSpace(rows[i][8]), Valid: rows[i][8] != ""}.String
		}
		//registered_residence_type tinyint COMMENT '户籍情况:0.未知1.东湖户籍2.北京市其他户籍3.外地户籍',
		if len(rows[i]) > 9 && rows[i][9] != "" {
			num, err := strconv.Atoi(rows[i][9])
			if err == nil {
				person.RegisteredResidenceType = num
			}
		}
		//registered_residence VARCHAR(200) COMMENT '户籍地',
		if len(rows[i]) > 10 {
			person.RegisteredResidence = sql.NullString{String: strings.TrimSpace(rows[i][10]), Valid: rows[i][10] != ""}.String
		}
		//telephone VARCHAR(50) COMMENT '联系方式',
		if len(rows[i]) > 11 {
			person.Telephone = sql.NullString{String: strings.TrimSpace(rows[i][11]), Valid: rows[i][11] != ""}.String
		}
		//first_contact VARCHAR(50) COMMENT '第一联系人',
		if len(rows[i]) > 12 {
			person.FirstContact = sql.NullString{String: strings.TrimSpace(rows[i][12]), Valid: rows[i][12] != ""}.String
		}
		//elder_relationship VARCHAR(50) COMMENT '关系(老人之)',
		if len(rows[i]) > 13 {
			person.ElderRelationship = sql.NullString{String: strings.TrimSpace(rows[i][13]), Valid: rows[i][13] != ""}.String
		}
		//elder_contact_phone VARCHAR(20) COMMENT '联系电话',
		if len(rows[i]) > 14 {
			person.ElderContactPhone = sql.NullString{String: strings.TrimSpace(rows[i][14]), Valid: rows[i][14] != ""}.String
		}
		//disability_level VARCHAR(100) COMMENT '失能等级',
		if len(rows[i]) > 15 {
			person.DisabilityLevel = sql.NullString{String: strings.TrimSpace(rows[i][15]), Valid: rows[i][15] != ""}.String
		}
		//is_low_income tinyint COMMENT '是否低保2否，1是',
		if len(rows[i]) > 16 {
			num := strings.TrimSpace(rows[i][16])
			if num == "是" {
				person.IsLowIncome = 1
			} else if num == "否" {
				person.IsLowIncome = 2
			}
		}
		//is_low_income2 tinyint COMMENT '是否低收入2否，1是',
		if len(rows[i]) > 17 {
			num := strings.TrimSpace(rows[i][17])
			if num == "是" {
				person.IsLowIncome2 = 1
			} else if num == "否" {
				person.IsLowIncome2 = 2
			}
		}
		//is_destitute tinyint COMMENT '是否特困供养2否，1是',
		if len(rows[i]) > 18 {
			num := strings.TrimSpace(rows[i][18])
			if num == "是" {
				person.IsDestitute = 1
			} else if num == "否" {
				person.IsDestitute = 2
			}
		}
		//is_family_planning_special tinyint COMMENT '是否计划生育特殊家庭2否，1是',
		if len(rows[i]) > 19 {
			num := strings.TrimSpace(rows[i][19])
			if num == "是" {
				person.IsFamilyPlanningSpecial = 1
			} else if num == "否" {
				person.IsFamilyPlanningSpecial = 2
			}
		}
		//disability_category VARCHAR(100) COMMENT '残疾类别及等级',
		if len(rows[i]) > 20 {
			person.DisabilityCategory = sql.NullString{String: strings.TrimSpace(rows[i][20]), Valid: rows[i][20] != ""}.String
		}
		//is_living_alone tinyint COMMENT '是否独居2否1是',
		if len(rows[i]) > 21 {
			num := strings.TrimSpace(rows[i][21])
			if num == "是" {
				person.IsLivingAlone = 1
			} else if num == "否" {
				person.IsLivingAlone = 2
			}
		}
		//is_empty_nest tinyint COMMENT '是否空巢2否1是',
		if len(rows[i]) > 22 {
			num := strings.TrimSpace(rows[i][22])
			if num == "是" {
				person.IsEmptyNest = 1
			} else if num == "否" {
				person.IsEmptyNest = 2
			}
		}
		//is_orphaned tinyint COMMENT '是否孤寡2否1是',
		if len(rows[i]) > 23 {
			num := strings.TrimSpace(rows[i][23])
			if num == "是" {
				person.IsOrphaned = 1
			} else if num == "否" {
				person.IsOrphaned = 2
			}
		}
		//other_situation TEXT COMMENT '其它情况',
		if len(rows[i]) > 24 {
			person.OtherSituation = sql.NullString{String: strings.TrimSpace(rows[i][24]), Valid: rows[i][24] != ""}.String
		}
		//is_needs_focus tinyint COMMENT '是否需重点关注2否1是',
		if len(rows[i]) > 25 {
			num := strings.TrimSpace(rows[i][25])
			if num == "是" {
				person.IsNeedsFocus = 1
			} else if num == "否" {
				person.IsNeedsFocus = 2
			}
		}
		//is_in_group tinyint COMMENT '是否入群2否1是',
		if len(rows[i]) > 28 {
			num := strings.TrimSpace(rows[i][28])
			if num == "是" {
				person.IsInGroup = 1
			} else if num == "否" {
				person.IsInGroup = 2
			}
		}
		//is_private_message tinyint COMMENT '是否私信2否1是',
		if len(rows[i]) > 29 {
			num := strings.TrimSpace(rows[i][29])
			if num == "是" {
				person.IsPrivateMessage = 1
			} else if num == "否" {
				person.IsPrivateMessage = 2
			}
		}
		//has_pet tinyint COMMENT '家有宠物2否1是',
		if len(rows[i]) > 30 {
			num := strings.TrimSpace(rows[i][30])
			if num == "是" {
				person.HasPet = 1
			} else if num == "否" {
				person.HasPet = 2
			}
		}
		//last_contact_time VARCHAR(50) COMMENT '末次联系时间',
		if len(rows[i]) > 31 {
			lastTime := strings.TrimSpace(rows[i][31])
			person.LastContactTime = lastTime
		}
		//other_info TEXT COMMENT '其他',
		if len(rows[i]) > 32 {
			person.OtherInfo = strings.TrimSpace(rows[i][32])
		}
		if person.BuildingNumber == "" {
			fmt.Println(person)
		}

		persons = append(persons, person)
	}
	return persons, nil
}

// savePersonInfo 将人员信息保存到数据库中
func (s *UpdateExcelDataService) savePersonInfo(person models.Person) error {
	return s.db.Create(&person).Error
}

func (s *UpdateExcelDataService) savePersons(persons []models.Person) error {
	return s.db.CreateInBatches(persons, 500).Error
}

/**
 * 从字符串中提取数字  从sheet名字中提取楼号
 * @param str 输入的字符串
 * @return 返回字符串中的第一个数字，如果没有找到则返回空字符串
 */
func extractNumbers(str string) string {
	// 编译正则表达式：匹配一个或多个数字
	re := regexp.MustCompile(`\d+`)

	// 查找第一个匹配的数字
	matches := re.FindStringSubmatch(str)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// parseBuildingAddress 函数用于解析地址字符串，提取出楼栋、单元和房间号信息
// 处理117楼特殊地址类型
// 参数:
//
//	addr: 待解析的地址字符串
//
// 返回值:
//
//	building: 楼栋号
//	unit: 单元号
//	room: 房间号
func parseBuildingAddress(addr string) (building, unit, room string) {
	// 正则表达式匹配：数字+"楼" + 数字+"单元" + 数字
	re := regexp.MustCompile(`(\d+)楼(\d+)单元(\d+)`)
	matches := re.FindStringSubmatch(addr)

	if len(matches) == 4 {
		return matches[1], matches[2], matches[3]
	}

	// 如果没有匹配到，尝试更宽松的正则
	re2 := regexp.MustCompile(`(\d+)[楼幢]?(\d*)[单元]?(\d*)`)
	matches2 := re2.FindStringSubmatch(addr)

	if len(matches2) >= 2 {
		building = matches2[1]
		if len(matches2) >= 3 && matches2[2] != "" {
			unit = matches2[2]
		} else {
			unit = "1" // 默认1单元
		}
		if len(matches2) >= 4 && matches2[3] != "" {
			room = matches2[3]
		} else {
			room = "未知"
		}
	}

	return building, unit, room
}
