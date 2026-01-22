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
	handler := handlers.NewUpdateExecDataHandler(db)
	handler.UpdateExecData("./web/alldata.xlsx")

	//num, err := strconv.Atoi("北京户籍")
	//if err == nil {
	//	fmt.Println(num)
	//}
}
