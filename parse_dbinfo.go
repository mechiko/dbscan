package dbscan

import (
	"fmt"
	"path/filepath"
)

// name driver должны быть заполнены
// file если пусто будет вычислен из name добавлением .db
func (d *dbs) ParseDbInfo(dbPath string, info *DbInfo) (dbi *DbInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if !filepath.IsAbs(info.File) {
		info.File = filepath.Join(dbPath, info.File)
	}
	dbi = &DbInfo{
		File:   info.File,
		Driver: info.Driver,
		Name:   info.Name,
		Host:   info.Host,
		User:   info.User,
		Pass:   info.Pass,
		Exists: false,
	}
	if err := d.IsConnected(dbi); err != nil {
		return nil, fmt.Errorf("is connected %w", err)
	} else {
		dbi.Exists = true
	}

	return dbi, nil
}
