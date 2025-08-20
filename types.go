package dbscan

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

const modError = "dbscan"

type Apper interface {
	Logger() *zap.SugaredLogger
	Pwd() string
	ConfigPath() string
	DbPath() string
}

type dbs struct {
	Apper
	path  string
	infos ListDbInfoForScan
}

// dbPath путь сканирования БД A3 и 4Z
// listInfo описатели для всех БД
func New(app Apper, listInfo ListDbInfoForScan, dbPath string) (d *dbs, err error) {
	d = &dbs{
		Apper: app,
		path:  dbPath,
		infos: make(ListDbInfoForScan),
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("panic %s: %v", modError, r))
		}
	}()

	fsrarId := ""
	file4z := ""
	dbType := "sqlite"
	// если Config есть в списке
	if config := listInfo.Info(Config); config != nil {
		if config.Path == "" {
			config.Path = dbPath
		}
		configParsedInfo, err := d.ParseDbInfo(config)
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
	if other := listInfo.Info(Other); other != nil {
		if other.Path == "" {
			other.Path = dbPath
		}
		otherParsedInfo, err := d.ParseDbInfo(other)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse other info %w", err)
		}
		d.infos[Other] = otherParsedInfo
	}
	if a3 := listInfo.Info(A3); a3 != nil {
		if a3.Path == "" {
			a3.Path = dbPath
		}
		if a3.Driver == "" {
			a3.Driver = dbType
		}
		if a3.Name == "" {
			a3.Name = fsrarId
		}
		if a3.File == "" {
			a3.File = a3.Name + ".db"
		}
		a3ParsedInfo, err := d.ParseDbInfo(a3)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse a3 info %w", err)
		}
		d.infos[A3] = a3ParsedInfo
	}
	if trueZnak := listInfo.Info(TrueZnak); trueZnak != nil {
		if trueZnak.Path == "" {
			trueZnak.Path = dbPath
		}
		if trueZnak.Driver == "" {
			trueZnak.Driver = dbType
		}
		if trueZnak.Name == "" {
			trueZnak.Name = file4z
		}
		if trueZnak.File == "" {
			trueZnak.File = trueZnak.Name + ".db"
		}
		trueParsedInfo, err := d.ParseDbInfo(trueZnak)
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
