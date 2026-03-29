package mageutil

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	FDIMRoot               string
	FDIMOutputConfig       string
	FDIMOutput             string
	FDIMOutputTools        string
	FDIMOutputTmp          string
	FDIMOutputLogs         string
	FDIMOutputBin          string
	FDIMOutputBinPath      string
	FDIMOutputBinToolPath  string
	FDIMInitErrLogFile     string
	FDIMInitLogFile        string
	FDIMOutputHostBin      string
	FDIMOutputHostBinTools string
)

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error getting current directory: " + err.Error())
	}

	FDIMRoot = currentDir

	FDIMOutputConfig = filepath.Join(FDIMRoot, "config") + string(filepath.Separator)
	FDIMOutput = filepath.Join(FDIMRoot, "_output") + string(filepath.Separator)

	FDIMOutputTools = filepath.Join(FDIMOutput, "tools") + string(filepath.Separator)
	FDIMOutputTmp = filepath.Join(FDIMOutput, "tmp") + string(filepath.Separator)
	FDIMOutputLogs = filepath.Join(FDIMOutput, "logs") + string(filepath.Separator)
	FDIMOutputBin = filepath.Join(FDIMOutput, "bin") + string(filepath.Separator)

	FDIMOutputBinPath = filepath.Join(FDIMOutputBin, "platforms") + string(filepath.Separator)
	FDIMOutputBinToolPath = filepath.Join(FDIMOutputBin, "tools") + string(filepath.Separator)

	FDIMInitErrLogFile = filepath.Join(FDIMOutputLogs, "FDIM-init-err.log")
	FDIMInitLogFile = filepath.Join(FDIMOutputLogs, "FDIM-init.log")

	FDIMOutputHostBin = filepath.Join(FDIMOutputBinPath, OsArch()) + string(filepath.Separator)
	FDIMOutputHostBinTools = filepath.Join(FDIMOutputBinToolPath, OsArch()) + string(filepath.Separator)

	dirs := []string{
		FDIMOutputConfig,
		FDIMOutput,
		FDIMOutputTools,
		FDIMOutputTmp,
		FDIMOutputLogs,
		FDIMOutputBin,
		FDIMOutputBinPath,
		FDIMOutputBinToolPath,
		FDIMOutputHostBin,
		FDIMOutputHostBinTools,
	}

	for _, dir := range dirs {
		createDirIfNotExist(dir)
	}
}

func createDirIfNotExist(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", dir, err)
		os.Exit(1)
	}
}

// GetBinFullPath constructs and returns the full path for the given binary name.
func GetBinFullPath(binName string) string {
	binFullPath := filepath.Join(FDIMOutputHostBin, binName)
	return binFullPath
}

// GetToolFullPath constructs and returns the full path for the given tool name.
func GetToolFullPath(toolName string) string {
	toolFullPath := filepath.Join(FDIMOutputHostBinTools, toolName)
	return toolFullPath
}
