package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"PLMS/internal/models"
	"PLMS/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PersonHandler struct {
	db      *gorm.DB
	service *services.PersonService
}

func NewPersonHandler(db *gorm.DB) *PersonHandler {
	return &PersonHandler{
		db:      db,
		service: services.NewPersonService(db),
	}
}

func (p *PersonHandler) GetPersons(c *gin.Context) {
	var filter models.PersonFilter

	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	persons, total, err := p.service.GetPersons(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    persons,
		"total":   total,
		"current": filter.Page,
	})
}

func (p *PersonHandler) GetRooms(c *gin.Context) {
	var filter models.PersonFilter

	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	roomList, total, err := p.service.GetRooms(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    roomList,
		"total":   total,
		"current": filter.Page,
	})
}

func (p *PersonHandler) GetPersonInfo(c *gin.Context) {
	personId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	personInfo, err := p.service.GetPersonInfo(personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": personInfo,
	})
}
func (p *PersonHandler) GetPersonInfoByRoom(c *gin.Context) {
	buildingNumber := c.Query("buildingNumber")
	unitNumber := c.Query("unitNumber")
	roomNumber := c.Query("roomNumber")
	personInfos, err := p.service.GetPersonInfoByRoom(buildingNumber, unitNumber, roomNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": personInfos,
	})
}

func (p *PersonHandler) GetBuildingNumbers(c *gin.Context) {
	buildingNumbers, err := p.service.GetBuildingNumbers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": buildingNumbers,
	})
}

func (p *PersonHandler) GetUnitNumbersByBuildingNumber(c *gin.Context) {
	buildingNumber := c.Query("buildingNumber")
	unitNumbers, err := p.service.GetUnitNumbersByBuildingNumber(buildingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": unitNumbers,
	})

}

func (p *PersonHandler) GetPersonStatistics(c *gin.Context) {
	buildingNumber := c.Query("buildingNumber")
	stat, err := p.service.GetPersonStatistics(buildingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": stat,
	})
}

// ExportPersons 导出人员信息到 Excel
// POST /api/v1/exportPersons
func (p *PersonHandler) ExportPersons(c *gin.Context) {
	var filter models.PersonFilter

	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 查询数据（不分页）
	persons, err := p.service.ExportPersons(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询失败",
		})
		return
	}

	if len(persons) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "没有数据可导出",
		})
		return
	}

	// 如果没有指定导出字段，使用默认字段
	if len(filter.ShowFields) == 0 {
		filter.ShowFields = []string{
			"building_number", "unit_number", "room_number",
			"name", "id_card", "age", "gender", "telephone",
		}
	}

	// 生成 Excel
	f := excelize.NewFile()
	sheetName := "人员信息"
	f.SetSheetName("Sheet1", sheetName)

	// 写入表头
	for colIdx, field := range filter.ShowFields {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, models.GetExportFieldHeader(field))
	}

	// 写入数据
	for rowIdx, person := range persons {
		for colIdx, field := range filter.ShowFields {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, person.GetExportValue(field))
		}
	}

	// 设置列宽
	for colIdx := range filter.ShowFields {
		col, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// 设置响应头
	filename := fmt.Sprintf("人员信息_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")

	// 写入响应
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "文件生成失败",
		})
	}
}

// GetExportFields 获取可导出的字段列表
// GET /api/v1/exportFields
func (p *PersonHandler) GetExportFields(c *gin.Context) {
	fields := []map[string]string{
		// 基本信息
		{"field": "building_number", "header": models.GetExportFieldHeader("building_number")},
		{"field": "unit_number", "header": models.GetExportFieldHeader("unit_number")},
		{"field": "room_number", "header": models.GetExportFieldHeader("room_number")},
		{"field": "name", "header": models.GetExportFieldHeader("name")},
		{"field": "id_card", "header": models.GetExportFieldHeader("id_card")},
		{"field": "age", "header": models.GetExportFieldHeader("age")},
		{"field": "gender", "header": models.GetExportFieldHeader("gender")},
		// 居住信息
		{"field": "is_permanent", "header": models.GetExportFieldHeader("is_permanent")},
		{"field": "housing_situation", "header": models.GetExportFieldHeader("housing_situation")},
		{"field": "property_nature", "header": models.GetExportFieldHeader("property_nature")},
		{"field": "registered_residence_type", "header": models.GetExportFieldHeader("registered_residence_type")},
		{"field": "registered_residence", "header": models.GetExportFieldHeader("registered_residence")},
		// 联系方式
		{"field": "telephone", "header": models.GetExportFieldHeader("telephone")},
		{"field": "first_contact", "header": models.GetExportFieldHeader("first_contact")},
		{"field": "elder_relationship", "header": models.GetExportFieldHeader("elder_relationship")},
		{"field": "elder_contact_phone", "header": models.GetExportFieldHeader("elder_contact_phone")},
		// 特殊情况
		{"field": "disability_level", "header": models.GetExportFieldHeader("disability_level")},
		{"field": "is_low_income", "header": models.GetExportFieldHeader("is_low_income")},
		{"field": "is_low_income2", "header": models.GetExportFieldHeader("is_low_income2")},
		{"field": "is_destitute", "header": models.GetExportFieldHeader("is_destitute")},
		{"field": "is_family_planning_special", "header": models.GetExportFieldHeader("is_family_planning_special")},
		{"field": "disability_category", "header": models.GetExportFieldHeader("disability_category")},
		{"field": "is_living_alone", "header": models.GetExportFieldHeader("is_living_alone")},
		{"field": "is_empty_nest", "header": models.GetExportFieldHeader("is_empty_nest")},
		{"field": "is_orphaned", "header": models.GetExportFieldHeader("is_orphaned")},
		{"field": "special_situation", "header": models.GetExportFieldHeader("special_situation")},
		{"field": "is_needs_focus", "header": models.GetExportFieldHeader("is_needs_focus")},
		// 电动车信息
		{"field": "has_electric_car", "header": models.GetExportFieldHeader("has_electric_car")},
		{"field": "license_plate", "header": models.GetExportFieldHeader("license_plate")},
		{"field": "brand_model", "header": models.GetExportFieldHeader("brand_model")},
		// 其他信息
		{"field": "is_in_group", "header": models.GetExportFieldHeader("is_in_group")},
		{"field": "is_private_message", "header": models.GetExportFieldHeader("is_private_message")},
		{"field": "has_pet", "header": models.GetExportFieldHeader("has_pet")},
		{"field": "last_contact_time", "header": models.GetExportFieldHeader("last_contact_time")},
		{"field": "other_info", "header": models.GetExportFieldHeader("other_info")},
		// 党员信息（如有）
		{"field": "is_cp", "header": models.GetExportFieldHeader("is_cp")},
		{"field": "nationality", "header": models.GetExportFieldHeader("nationality")},
		{"field": "education", "header": models.GetExportFieldHeader("education")},
		{"field": "cp_joining_day", "header": models.GetExportFieldHeader("cp_joining_day")},
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fields,
	})
}
