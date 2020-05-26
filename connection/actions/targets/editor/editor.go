package editor

import (
	"fmt"
	"fyssl/config"
	"fyssl/connection/actions/targets/base"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"slogger"
	"strings"
	"time"
)

type Editor struct {
	packetPath string
	editorProgram string
	editorArgs string
}

const (
	fileArg = "{file}"
	tempFileFormat = "fyssl-temp-file%d"
	maxFileNameNumber = 100000
	fileChangeErrorMessage = "Error while waiting for file change in file-%s, error:%+v"
	fileCantBeReadErrorMessage = "Error reading file-%s, sending received data"
)

func (e *Editor) ProcessTarget(buffer *[]byte) (*[]byte, error) {
	e.createPacketFile(buffer)
	e.createArgs()
	c := exec.Command(e.editorProgram, e.editorArgs)
	slogger.Debug(e.editorProgram)
	slogger.Debug(e.editorArgs)
	c.Start()
	e.WaitForFileChange()
	fileData := e.getFileData(buffer)
	slogger.Info(fmt.Sprintf("Getting data from %s and deleting the file!", e.packetPath))
	e.deleteFile()
	return fileData, nil
}

func (e *Editor) deleteFile() {
	err := os.Remove(e.packetPath)
	if err != nil {
		slogger.Error(fmt.Sprintf("Error deleting packet data file-%s,%+v", e.packetPath, err))
	}
}

func (e *Editor) createPacketFile(buffer *[]byte) error {
	e.packetPath = e.getTempFilePath()
	err := ioutil.WriteFile(e.packetPath, *buffer, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (e *Editor) getFileData(buffer *[]byte) *[]byte {
	if _, err := os.Stat(e.packetPath); err == nil {
		content, err := ioutil.ReadFile(e.packetPath)     // the file is inside the local directory
		if err != nil {
			slogger.Error(fmt.Sprintf(fileCantBeReadErrorMessage, e.packetPath))
			return buffer
		}
		return &content
	} else {
		slogger.Error(fmt.Sprintf(fileCantBeReadErrorMessage, e.packetPath))
		return buffer
	}
}

func (e *Editor) WaitForFileChange() {
	slogger.Info(fmt.Sprintf("Waiting for modification of file-%s", e.packetPath))
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slogger.Error(fmt.Sprintf(fileChangeErrorMessage, e.packetPath, err))
		watcher.Close()
		return
	}

	err = watcher.Add(e.packetPath)
	if err != nil {
		slogger.Error(fmt.Sprintf(fileChangeErrorMessage, e.packetPath, err))
		watcher.Close()
		return
	}

	select {
	case _ = <-watcher.Events:
		watcher.Close()
		return
	case err := <-watcher.Errors:
		slogger.Error(fmt.Sprintf(fileChangeErrorMessage, e.packetPath, err))
		watcher.Close()
		return
	}
	watcher.Close()
}


func (e *Editor) createArgs() {
	e.editorArgs = strings.Replace(e.editorArgs, fileArg, e.packetPath, -1)
}

func (e *Editor) getTempFilePath() string {
	for {
		fileName := e.getTempFileName()
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return fileName
		}
	}
}

func (e *Editor) getTempFileName() string {
	cfg := config.GetConfig()
	return path.Join(cfg.LogPath, fmt.Sprintf(tempFileFormat, e.getRandomNumber(maxFileNameNumber)))
}

func (e *Editor) getRandomNumber(maxRange int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(maxRange)
}

func (e *Editor) Close() {}

func CreateEditor(sock net.Conn, connection *config.Connection, action *config.Action) base.Target {
	dict, _ := action.TargetParams.(map[string]interface{})
	e := Editor{
		editorProgram: dict["editor"].(string),
		editorArgs: dict["args"].(string),
	}
	return &e
}

