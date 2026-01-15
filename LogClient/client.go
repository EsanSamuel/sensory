package logClient

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	conn    net.Conn
	Project string
	Service string
}

type LogEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Project   string `json:"project"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}

func New(project, service, addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Warning: Failed to connect to log server: %v", err)
		return nil, err
	}
	return &Client{conn: conn, Project: project, Service: service}, nil
}

func (c *Client) Send(level string, msg string) error {
	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   c.Project,
		Service:   c.Service,
		Message:   msg,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Println(err)
	}

	data = append(data, '\n')
	_, err = c.conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
