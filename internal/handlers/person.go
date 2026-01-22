package handlers

import (
	"PLMS/internal/models"
	"PLMS/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
		"data":  persons,
		"total": total,
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
		"data":  roomList,
		"total": total,
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
