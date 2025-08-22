package main

import (
	"os"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
	"go.uber.org/zap"
)

// .nevakod\4zupper

func main() {
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
	// 	os.Exit(1)
	// }
	// dir, _ := os.Getwd()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // Flushes buffer, if any
	sugar := logger.Sugar()
	list := make(dbscan.ListDbInfoForScan)
	// list[dbscan.Config] = &dbscan.DbInfo{
	// 	File:   "config.db",
	// 	Driver: "sqlite",
	// }
	// list[dbscan.Other] = &dbscan.DbInfo{
	// 	File:   "4zupper.db",
	// 	Name:   "zupper",
	// 	Driver: "sqlite",
	// 	Path:   `.nevakod\4zupper`,
	// }
	list[dbscan.A3] = &dbscan.DbInfo{}
	// list[dbscan.TrueZnak] = &dbscan.DbInfo{}
	dbs, err := dbscan.New(list, "")
	if err != nil {
		sugar.Errorf("ошибка создания %v", err)
		utility.MessageBox("ошибка", err.Error())
		os.Exit(-1)
	}
	sugar.Info(dbs)
}
