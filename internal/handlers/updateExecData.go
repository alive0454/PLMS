package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"PLMS/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateExecDataHandler struct {
	db      *gorm.DB
	service *services.UpdateExcelDataService
}

func NewUpdateExecDataHandler(db *gorm.DB) *UpdateExecDataHandler {
	return &UpdateExecDataHandler{
		db:      db,
		service: services.NewUpdateExcelDataService(db),
	}
}

func (h *UpdateExecDataHandler) UpdateExecData(filePath string) (*services.ImportResult, error) {
	return h.service.ImportExcelData(filePath)
}

// ImportExcel 导入Excel文件（HTTP接口）
// POST /api/v1/import/excel
func (h *UpdateExecDataHandler) ImportExcel(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 检查文件类型
	ext := filepath.Ext(header.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "只支持 .xlsx 或 .xls 格式的Excel文件",
		})
		return
	}

	// 创建临时目录
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建临时目录失败: " + err.Error(),
		})
		return
	}

	// 保存文件到临时目录
	tempFilePath := filepath.Join(tempDir, header.Filename)
	if err := c.SaveUploadedFile(header, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存文件失败: " + err.Error(),
		})
		return
	}

	// 处理完成后删除临时文件
	defer os.Remove(tempFilePath)

	// 调用导入服务
	result, err := h.service.ImportExcelData(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导入失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Excel导入成功",
		"data": gin.H{
			"filename":     header.Filename,
			"size":         header.Size,
			"totalSheets":  result.TotalSheets,
			"totalPersons": result.TotalPersons,
			"details":      result.Details,
		},
	})
}
