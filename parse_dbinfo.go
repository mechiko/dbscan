package dbscan

import (
	"fmt"
	"path/filepath"
)

// name driver должны быть заполнены
// file если пусто будет вычислен из name добавлением .db
func (d *dbs) ParseDbInfo(info *DbInfo) (dbi *DbInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	// Derive effective file path without mutating input
	filePath := info.File
	if filePath == "" && info.Driver == "sqlite" {
		if info.Name == "" {
			return nil, fmt.Errorf("%s отсутствует имя базы данных для sqlite", modError)
		}
		filePath = filepath.Join(info.Path, info.Name+".db")
	}
	if filePath != "" && !filepath.IsAbs(filePath) {
		filePath = filepath.Join(info.Path, filePath)
	}
	dbi = &DbInfo{
		File:   filePath,
		Driver: info.Driver,
		Name:   info.Name,
		Host:   info.Host,
		User:   info.User,
		Pass:   info.Pass,
		Exists: false,
	}
	if err := d.IsConnected(dbi); err != nil {
		return nil, fmt.Errorf("ошибка %w", err)
	} else {
		dbi.Exists = true
	}

	return dbi, nil
}
