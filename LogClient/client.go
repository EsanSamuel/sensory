package logClient

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/EsanSamuel/sensory/models"
)

type Client struct {
	conn      net.Conn
	Project   string
	Service   string
	ApiKey    string
	ProjectId string
	UserId    string
	noOp      bool // flag to mark dummy client
}

func New(apikey, projectName, addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		Project: projectName,
		ApiKey:  apikey,
		noOp:    false,
	}

	return client, nil
}

func NewNoOp() *Client {
	return &Client{noOp: true}
}

func getLocation() (file string, line int, function string) {
	pc, file, line, ok := runtime.Caller(3)

	if !ok {
		return "unknown", 0, "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "unknown"
	}
	return file, line, fn.Name()
}

func (c *Client) Send(level, msg string) error {
	if c == nil || c.noOp {
		fmt.Println("[noop logger]", level, msg)
		return nil
	}

	file, line, fn := getLocation()

	entry := models.LogEntry{
		Level:     level,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   c.Project,
		Service:   c.Service,
		Message:   msg,
		ApiKey:    c.ApiKey,
		Runtime: models.Runtime{
			File: file,
			Line: line,
			Fn:   fn,
		},
	}

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data = append(data, '\n')
	if c.conn != nil {
		_, err = c.conn.Write(data)
		if err != nil {
			fmt.Println("TCP write failed:", err)
		}
	}
	return err
}
