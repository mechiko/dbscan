package dbscan

import (
	"errors"
	"fmt"
	"os"

	"github.com/mechiko/utility"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mssql"
	"github.com/upper/db/v4/adapter/sqlite"
)

func (d *dbs) IsConnected(info *DbInfo) (err error) {
	var dbSess db.Session
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	switch info.Driver {
	case "sqlite":
		if info.File == "" {
			return fmt.Errorf("%s отсутствуют имя файла базы данных для sqlite", modError)
		}
		if !utility.PathOrFileExists(info.File) {
			return fmt.Errorf("%s отсутствует файл базы данных %s для sqlite", modError, info.File)
		}
		// если указан не файл а путь к каталогу
		if st, statErr := os.Stat(info.File); statErr != nil || !st.Mode().IsRegular() {
			return fmt.Errorf("%s путь %s не является файлом sqlite", modError, info.File)
		}
		uri := info.SqliteUri()
		dbSess, err = sqlite.Open(uri)
		if err != nil {
			return fmt.Errorf("%s ошибка подключения %v", modError, err)
		}
		defer func() {
			if errSess := dbSess.Close(); errSess != nil {
				// Go 1.20+: joins parse error (if any) with close error
				err = errors.Join(err, fmt.Errorf("close session %w", errSess))
			}
		}()
	case "mssql":
		if info.Name == "" {
			return fmt.Errorf("%s отсутствуют имя базы данных для Other", modError)
		}
		uri := info.MssqlUri()
		dbSess, err = mssql.Open(uri)
		if err != nil {
			return fmt.Errorf("%s %s", modError, err.Error())
		}
		defer func() {
			if errSess := dbSess.Close(); errSess != nil {
				// Go 1.20+: joins parse error (if any) with close error
				err = errors.Join(err, fmt.Errorf("close session %w", errSess))
			}
		}()
	default:
		return fmt.Errorf("%s ошибка driver %v", modError, info.Driver)
	}
	err = dbSess.Ping()
	if err != nil {
		return fmt.Errorf("%s ошибка ping %v", modError, err)
	}
	// пинг успешен
	return nil

}
