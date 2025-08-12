package dbscan

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mechiko/utility"
)

var reA3DbName = regexp.MustCompile(`^0[0-9]{11}\.db$`)

func findA3DbName() string {
	if files, err := utility.FilteredSearchOfDirectoryTree(reA3DbName, ""); err != nil {
		return ""
	} else {
		if len(files) == 0 {
			return ""
		}
		return files[0]
	}
}

func findA3Name() string {
	findName := findA3DbName()
	if findName == "" {
		return ""
	}
	ext := filepath.Ext(findName)
	return strings.TrimSuffix(findName, ext)
}
