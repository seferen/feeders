package file

import (
	"bufio"
	"os"

	"go.uber.org/zap"
)

type file struct {
	Log        *zap.Logger
	sourceFile *os.File
	sourcePath string
	scanner    *bufio.Scanner
}

func GetFeeder(path string) *file {
	log, _ := zap.NewProduction()

	result := &file{Log: log, sourcePath: path}

	result.openFile()

	return result
}

func (f *file) openFile() {
	var err error

	f.sourceFile, err = os.Open(f.sourcePath)

	if err != nil {
		f.Log.Fatal("Open file", zap.String("file name", f.sourcePath), zap.String("error", err.Error()))
	}

	f.Log.Info("open file", zap.String("file name", f.sourceFile.Name()))

	f.scanner = bufio.NewScanner(f.sourceFile)
}

func (f *file) Feed() string {

	if f.scanner.Scan() {
		text := f.scanner.Text()

		f.Log.Debug("Read string", zap.String("string line", text))

		return text
	}

	f.openFile()

	return f.Feed()

}

func (f *file) Close() {
	f.sourceFile.Close()

}
