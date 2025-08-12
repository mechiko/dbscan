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
	driver string
	path   string
	// self           *DbInfo
	// a3             *DbInfo
	// znak           *DbInfo
	// config         *DbInfo
	// configFileName string // config.db алкохелпа
	infos ListDbInfoForScan
}

// dbPath путь сканирования БД
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
		configParsedInfo, err := d.ParseDbInfo(app.DbPath(), config)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse config info %v", err)
		}
		d.infos[Config] = configParsedInfo
		file4z, err = d.fromConfig(config, "oms_id")
		if err != nil {
			return nil, fmt.Errorf("dbscan fromConfig %v", err)
		}
		fsrarId, err = d.fromConfig(config, "fsrar_id")
		if err != nil {
			return nil, fmt.Errorf("dbscan fromConfig %v", err)
		}
		dbType, err = d.fromConfig(config, "db_type")
		if err != nil {
			return nil, fmt.Errorf("dbscan fromConfig %v", err)
		}
		dbType = strings.ToLower(dbType)
	}
	if fsrarId == "" {
		fsrarId = findA3Name()
	}
	if other := listInfo.Info(Other); other != nil {
		otherParsedInfo, err := d.ParseDbInfo(app.ConfigPath(), other)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse other info %v", err)
		}
		d.infos[Other] = otherParsedInfo
	}
	if a3 := listInfo.Info(A3); a3 != nil {
		if a3.Driver == "" {
			a3.Driver = dbType
		}
		if a3.Name == "" {
			a3.Name = fsrarId
		}
		if a3.File == "" {
			a3.File = a3.Name + ".db"
		}
		a3ParsedInfo, err := d.ParseDbInfo(app.ConfigPath(), a3)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse a3 info %v", err)
		}
		d.infos[A3] = a3ParsedInfo
	}
	if trueZnak := listInfo.Info(TrueZnak); trueZnak != nil {
		if trueZnak.Driver == "" {
			trueZnak.Driver = dbType
		}
		if trueZnak.Name == "" {
			trueZnak.Name = file4z
		}
		if trueZnak.File == "" {
			trueZnak.File = trueZnak.Name + ".db"
		}
		trueParsedInfo, err := d.ParseDbInfo(app.ConfigPath(), trueZnak)
		if err != nil {
			return nil, fmt.Errorf("dbscan parse trueznak info %v", err)
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
