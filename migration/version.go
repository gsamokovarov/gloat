package migration

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Now           = time.Now()
	versionFormat = "20060319150405"
)

func generateVersion() int64 {
	version, _ := strconv.ParseInt(Now.Format(versionFormat), 10, 64)
	return version
}

func versionFromPath(path string) (string, error) {
	parts := strings.SplitN(filepath.Base(path), "_", 2)
	if len(parts) == 0 {
		return "", fmt.Errorf("cannot extract version from %s", path)
	}

	return parts[0], nil
}
