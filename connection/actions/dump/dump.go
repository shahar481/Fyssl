package dump

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/shahar481/fyssl/config"
	"net"
	"os"
	"path"
	"time"
)

const (
	logNameFormat = "dump-%s>%s"
	logEntryFormat = "%d-%d-%d %s>%s"
)

func Packet(buffer *[]byte, receiver net.Conn, sender net.Conn, connection *config.Connection, action *config.Action) {
	log := getLogPath(receiver, sender, action, connection)

	f, err := os.OpenFile(log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		glog.Errorf("Error opening the file for dumping in action:%s,%+v", action.Name, err)
	}

	defer f.Close()

	logEntry := createDumpLog(buffer, receiver, sender)
	if _, err = f.WriteString(logEntry); err != nil {
		glog.Errorf("Error dumping to file in action:%s,%+v", action.Name, err)
	}
}

func createDumpLog(buffer *[]byte, receiver net.Conn, sender net.Conn) string {
	time := time.Now()
	log := fmt.Sprintf(logEntryFormat, time.Year(), time.Month(), time.Day(), receiver.LocalAddr(), sender.LocalAddr())
	log = fmt.Sprintf("%s\n%s", log, hex.Dump(*buffer))
	return log
}

func createFolderPath(connection *config.Connection, action *config.Action) string {
	cfg := config.GetConfig()
	folderPath := path.Join(cfg.LogPath, connection.Name, action.Name)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		glog.Errorf("Couldn't create folders for %s-%+v", action.Name, err)
	}
	return folderPath
}

func getLogPath(receiver net.Conn, sender net.Conn, action *config.Action, connection *config.Connection) string {
	log, err := getCurrentLogFile(receiver, sender, action, connection)
	if err != nil {
		log, _ = getCurrentLogFile(sender, receiver, action, connection)
	}
	return log
}

func getCurrentLogFile(receiver net.Conn, sender net.Conn, action *config.Action, connection *config.Connection) (string, error) {
	log := path.Join(createFolderPath(connection, action), getLogName(receiver, sender))
	if _, err := os.Stat(log); !os.IsNotExist(err) {
		return log, nil
	}
	return log, errors.New("log file does not exist")
}

func getLogName(receiver net.Conn, sender net.Conn) string {
	return fmt.Sprintf(logNameFormat, receiver.LocalAddr(), sender.LocalAddr())
}