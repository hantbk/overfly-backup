package archive

import (
	"fmt"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
	"path"
	"path/filepath"
)

// Run archive
func Run(model config.ModelConfig) (err error) {
	logger.Info("----------- Archive Files ----------")

	tarPath := path.Join(model.DumpPath, "archive.tar")

	includes := model.Archive.GetStringSlice("includes")
	includes = cleanPaths(includes)

	if len(includes) == 0 {
		return fmt.Errorf("archive.includes have no config")
	}
	logger.Info("=> includes", len(includes), "rules")

	cmd := "tar -cPf " + tarPath

	excludes := model.Archive.GetStringSlice("excludes")
	excludes = cleanPaths(excludes)

	for _, exclude := range excludes {
		cmd += " --exclude='" + filepath.Clean(exclude) + "'"
	}

	helper.Exec(cmd, includes...)

	//
	logger.Info("----------- Archive Files ----------\n")
	return nil
}

func cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}
