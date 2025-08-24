package dbscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsTypes = []struct {
	name string
	err  bool
	list ListDbInfoForScan
	dir  string
}{
	// the table itself
	// 0 так мы не находим конфиг поэтому не находим и А3 будет ошибка по конфиг файлу
	{"test 0", true, ListDbInfoForScan{Config: &DbInfo{}, A3: &DbInfo{}}, ""},
	// 1 так мы находим конфиг и следом находим по нему mssql A3 но если sqlite и будет путь не указан не найдет саму БД
	{"test 1", false, ListDbInfoForScan{Config: &DbInfo{Path: "cmd"}, A3: &DbInfo{}}, ""},
	// 2 найдет обе БД
	{"test 2", false, ListDbInfoForScan{Config: &DbInfo{Path: "cmd"}, A3: &DbInfo{Driver: "sqlite", Path: "cmd"}}, ""},
	{"test 3", true, ListDbInfoForScan{Config: &DbInfo{File: "cmd/config.db", Driver: "sqlite"}}, ""},
	{"test 4", false, ListDbInfoForScan{Config: &DbInfo{File: "config.db", Driver: "sqlite"}}, "cmd"},
	{"test 5", false, ListDbInfoForScan{Other: &DbInfo{File: "4zupper.db", Driver: "sqlite"}}, "cmd"},
	{"test 6", false, ListDbInfoForScan{Other: &DbInfo{File: "4zupper.db", Driver: "sqlite"}}, "cmd/.nevakod/4zupper"},
}

func TestNew(t *testing.T) {
	// The execution loop
	for _, tt := range testsTypes {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.list, tt.dir)
			if tt.err {
				assert.NotNil(t, err)
			} else {
				// ожидаем отсутствие ошибки
				assert.Nil(t, err)
			}
		})
	}
}
