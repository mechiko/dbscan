package dbscan

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mechiko/utility"
)

const modError = "dbscan"

// type Apper interface {
// 	Logger() *zap.SugaredLogger
// 	Pwd() string
// 	ConfigPath() string
// 	DbPath() string
// }

type dbs struct {
	// Apper
	path  string
	infos ListDbInfoForScan
}

// dbPath путь сканирования БД A3 и 4Z
// listInfo описатели для всех БД
// пустой путь перезаписывается на текущий "."
func New(listInfo ListDbInfoForScan, dbPath string) (d *dbs, err error) {
	d = &dbs{
		// Apper: app,
		infos: make(ListDbInfoForScan),
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("panic %s: %v", modError, r))
		}
	}()
	if dbPath != "" && !utility.PathOrFileExists(dbPath) {
		return nil, fmt.Errorf("path %s not present", dbPath)
	}
	// для поиска файлов надо точку
	if dbPath == "" {
		dbPath = "."
	}
	d.path = dbPath
	fsrarId := ""
	file4z := ""
	dbType := "sqlite"
	// если Config есть в списке
	if config, ok := listInfo[Config]; ok && config != nil {
		if config.Path == "" {
			config.Path = dbPath
		}
		// only sqlite config.db
		config.Driver = "sqlite"
		config.File = "config.db"
		// проверяем структуру и пробуем коннект
		configParsedInfo, err := ParseDbInfo(config)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse config info %w", err)
		}
		d.infos[Config] = configParsedInfo
		file4z, err = d.fromConfig(config, "oms_id")
		if err != nil {
			return nil, fmt.Errorf("dbscan fromConfig %w", err)
		}
		fsrarId, err = d.fromConfig(config, "fsrar_id")
		if err != nil {
			return nil, fmt.Errorf("dbscan fromConfig %w", err)
		}
		dbType, err = d.fromConfig(config, "db_type")
		if err != nil {
			return nil, fmt.Errorf("dbscan fromConfig %w", err)
		}
		dbType = strings.ToLower(dbType)
	}
	if fsrarId == "" {
		fsrarId = findA3Name(dbPath)
	}
	if file4z == "" {
		file4z = find4zName(dbPath)
	}
	if other, ok := listInfo[Other]; ok && other != nil {
		if other.Path == "" {
			other.Path = dbPath
		}
		otherParsedInfo, err := ParseDbInfo(other)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse other info %w", err)
		}
		d.infos[Other] = otherParsedInfo
	}
	if a3, ok := listInfo[A3]; ok && a3 != nil {
		if a3.Path == "" {
			a3.Path = dbPath
		}
		if a3.Driver == "" {
			a3.Driver = dbType
		}
		if a3.Name == "" {
			a3.Name = fsrarId
		}
		a3ParsedInfo, err := ParseDbInfo(a3)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse a3 info %w", err)
		}
		d.infos[A3] = a3ParsedInfo
	}
	if trueZnak, ok := listInfo[TrueZnak]; ok && trueZnak != nil {
		if trueZnak.Path == "" {
			trueZnak.Path = dbPath
		}
		if trueZnak.Driver == "" {
			trueZnak.Driver = dbType
		}
		if trueZnak.Name == "" {
			trueZnak.Name = file4z
		}
		trueParsedInfo, err := ParseDbInfo(trueZnak)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse trueznak info %w", err)
		}
		d.infos[TrueZnak] = trueParsedInfo
	}
	return d, nil
}

func (d *dbs) Info(t DbInfoType) *DbInfo {
	if dbi, ok := d.infos[t]; ok {
		return dbi
	}
	return nil
}
