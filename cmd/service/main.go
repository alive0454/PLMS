package main

import (
	"PLMS/internal/config"
	"PLMS/internal/database"
	"PLMS/internal/handlers"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	//人口台账
	handler := handlers.NewUpdateExecDataHandler(db)
	result, err := handler.UpdateExecData("/Users/wangyao/GolandProjects/PLMS/web/alldata.xlsx")
	if err != nil {
		log.Fatal("导入失败:", err)
	}
	log.Printf("导入完成: %d 个工作表, %d 条数据", result.TotalSheets, result.TotalPersons)
	for _, detail := range result.Details {
		log.Println(detail)
	}
	//党员台账

}
