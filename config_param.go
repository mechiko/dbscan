package dbscan

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
)

type Parameters struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}

func (p *Parameters) Store(sess db.Session) db.Store {
	return sess.Collection("parameters")
}

func (dd *dbs) fromConfig(config *DbInfo, key string) (out string, err error) {
	var sess db.Session
	out = ""
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("panic %s: %v", config.File, r))
		}
	}()

	if config == nil {
		return "", fmt.Errorf("config dbinfo is nil")
	}
	switch config.Driver {
	case "sqlite":
		uri := sqlite.ConnectionURL{
			Database: filepath.Join(config.Path, config.File),
			Options: map[string]string{
				"mode":          "rw",
				"_journal_mode": "DELETE",
			},
		}
		sess, err = sqlite.Open(uri)
		if err != nil {
			return "", fmt.Errorf("dbscan:fromconfig %s", err.Error())
		}
		defer func() {
			if errClose := sess.Close(); err != nil {
				// Go 1.20+: joins parse error (if any) with close error
				err = errors.Join(err, fmt.Errorf("close %s: %w", config.File, errClose))
			}
		}()
	}
	param := &Parameters{}
	if err = sess.Get(param, db.Cond{"name": key}); err != nil {
		return "", fmt.Errorf("dbscan %s %v", config.File, err)
	}
	return param.Value, nil
}
