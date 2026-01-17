package logserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/EsanSamuel/sensory/helpers"
)

type Log struct {
	Level     string `json:"level"`
	TimeStamp string `json:"timestamp"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Project   string `json:"project"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
}

const serviceDir = "services"
const serviceLogFile = "service.log"
const errorLogFile = "error.log"
const infoLogFile = "info.log"
const fatalLogFile = "fatal.log"

func Initialize_Log() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer l.Close()

	fmt.Println("TCP Log Server listening on :9000")

	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		fmt.Println("Failed to create services folder:", err)
		return
	}

	ProcessLog()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Received log:", line)

		var entry Log
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Println(err)
			continue
		}

		logPath := filepath.Join(serviceDir, serviceLogFile)
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed to open log file:", err)
			continue
		}

		if _, err := f.WriteString(line + "\n"); err != nil {
			fmt.Println("Failed to write log:", err)
		}

		f.Close()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
	}
}

func ProcessLog() {
	logPath := filepath.Join(serviceDir, serviceLogFile)

	file, err := os.ReadFile(logPath)
	src, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No existing logs found.")
			return
		}
		fmt.Println("Error reading log file:", err)
		return
	}

	scanner := bufio.NewScanner(src)
	LogData := make(map[string][]Log)
	var entry Log

	for scanner.Scan() {
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Println(err)
			return
		}

		//Time := strings.Split(entry.TimeStamp, "T")
		//timeFormat := Time[0]
		//fmt.Println(timeFormat)

		data, err := os.ReadFile(".sensory.json")
		if err != nil {
			log.Println("err reading ", err)
		}

		var project helpers.ProjectMeta

		if err := json.Unmarshal(data, &project); err != nil {
			log.Println(err)
		}

		LogData[entry.Level] = append(LogData[entry.Level], Log{
			Project:   project.ProjectName,
			TimeStamp: entry.TimeStamp,
			Message:   entry.Message,
			Service:   entry.Service,
			ProjectID: project.ProjectId,
			UserID:    project.UserId,
		})

		switch entry.Level {
		case "INFO":
			logPath := filepath.Join(serviceDir, infoLogFile)
			writeLog(logPath, line)
		case "ERROR":
			logPath := filepath.Join(serviceDir, errorLogFile)
			writeLog(logPath, line)
		case "FATAL":
			logPath := filepath.Join(serviceDir, fatalLogFile)
			writeLog(logPath, line)

		}

	}

	fmt.Println("Error Logs: ", LogData["ERROR"])
	fmt.Println("Info Logs: ", LogData["INFO"])
	fmt.Println("Fatal Logs: ", LogData["FATAL"])

	fmt.Println("\n--- Existing Logs ---")
	fmt.Println(string(file))
	fmt.Printf("--- End of Logs ---\n")
}

func writeLog(logPath, line string) {
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	line = strings.TrimRight(line, "\n")
	if _, err := f.WriteString(line + "\n"); err != nil {
		fmt.Println("Failed to write log:", err)
	}

	f.Close()
}
