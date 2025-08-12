package dbscan

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mechiko/utility"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func find4zDbName(dir string) string {
	if files, err := utility.FilteredSearchOfDirectoryTree(uuidRegex, dir); err != nil {
		return ""
	} else {
		if len(files) == 0 {
			return ""
		}
		return files[0]
	}
}

func find4zName() string {
	findName := findA3DbName()
	if findName == "" {
		return ""
	}
	ext := filepath.Ext(findName)
	return strings.TrimSuffix(findName, ext)
}
