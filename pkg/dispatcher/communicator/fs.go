package communicator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
)

type fsCommunicator struct {
	inputPath  string
	outputPath string
}

func NewFSCommunicator(inputPath, outputPath string) (Communicator, error) {
	outputPathStat, err := os.Stat(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path \"%s\" does not exist", outputPath)
		}

		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied on accessing \"%s\"", outputPath)
		}

		return nil, err
	}

	if !outputPathStat.IsDir() {
		return nil, fmt.Errorf("path \"%s\" is not a directory", outputPath)
	}

	if err := os.MkdirAll(inputPath, 0755); err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied on accessing \"%s\"", inputPath)
		}

		return nil, err
	}

	return &fsCommunicator{
		inputPath:  inputPath,
		outputPath: outputPath,
	}, nil
}

func (c fsCommunicator) getFileContent(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (c fsCommunicator) removeFile(path string) error {
	return os.Remove(path)
}

func (c fsCommunicator) writeWriteContent(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0644)
}

func (c fsCommunicator) FetchData() ([]byte, error) {
	fileInfos, err := ioutil.ReadDir(c.outputPath)
	if err != nil {
		return []byte{}, err
	}

	if len(fileInfos) > 0 {
		// find the oldest file
		oldestFileInfo := fileInfos[0]
		for _, eachFileInfo := range fileInfos {
			if eachFileInfo.ModTime().Before(oldestFileInfo.ModTime()) {
				oldestFileInfo = eachFileInfo
			}
		}
		fileName := path.Join(c.outputPath, oldestFileInfo.Name())
		data, err := c.getFileContent(fileName)
		if err != nil {
			return []byte{}, err
		}
		if err := c.removeFile(fileName); err != nil {
			return []byte{}, err
		}

		return data, nil
	}

	// watch for new files
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return []byte{}, err
	}
	defer watcher.Close()

	outputFileNameChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go func(w *fsnotify.Watcher, d chan<- string, e chan<- error) {
		for {
			select {
			case event := <-w.Events:
				op := event.Op
				if op&(fsnotify.Create|fsnotify.Write) > 0 {
					d <- event.Name
					break
				}
			case err := <-w.Errors:
				e <- err
				break
			}
		}
	}(watcher, outputFileNameChan, errChan)

	if err := watcher.Add(c.outputPath); err != nil {
		return []byte{}, err
	}

	select {
	case outputFileName := <-outputFileNameChan:
		data, err := c.getFileContent(outputFileName)
		if err != nil {
			return []byte{}, err
		}
		if err := c.removeFile(outputFileName); err != nil {
			return []byte{}, err
		}
		return data, nil
	case err := <-errChan:
		return []byte{}, err
	}
}

func (c fsCommunicator) SendData(data []byte) ([]byte, error) {
	fileName := fmt.Sprintf("data-%d", time.Now().Unix())
	filePath := path.Join(c.inputPath, fileName)

	err := c.writeWriteContent(filePath, data)
	return []byte{}, err
}
