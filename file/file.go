package file

import (
	"bufio"
	"os"

	"go.uber.org/zap"
)

var files map[string]chan string = make(map[string]chan string)

type file struct {
	Log        *zap.Logger
	sourceFile *os.File
	sourcePath string
	scanner    *bufio.Scanner
	line       chan string
}

func GetFeeder(path string) *file {
	log, _ := zap.NewProduction()

	result := &file{Log: log, sourcePath: path}

	if quee, ok := files[result.sourcePath]; ok {
		result.line = quee
	} else {
		result.line = make(chan string)
		files[result.sourcePath] = result.line
	}

	result.openFile()

	go func() {
		for {
			switch result.scanner.Scan() {
			case true:
				result.line <- result.scanner.Text()
			case false:
				result.openFile()
				result.line <- result.scanner.Text()
			}
		}
	}()

	return result
}

func (f *file) openFile() {

	var err error

	f.sourceFile, err = os.Open(f.sourcePath)

	if err != nil {
		f.Log.Fatal("Open file", zap.String("file name", f.sourcePath), zap.String("error", err.Error()))
	}

	f.Log.Debug("open file", zap.String("file name", f.sourceFile.Name()))

	f.scanner = bufio.NewScanner(f.sourceFile)
}

func (f *file) Feed() string {

	return <-f.line

}

func (f *file) Close() {
	f.Log.Debug("call function close")
	f.sourceFile.Close()
	close(f.line)

}
