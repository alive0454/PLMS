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
	handler.UpdateExecData("./web/alldata.xlsx")
	//党员台账

}
