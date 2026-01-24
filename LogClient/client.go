package logClient

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/EsanSamuel/sensory/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	Project   string
	Service   string
	ApiKey    string
	ProjectId string
	UserId    string
	noOp      bool
}

func New(apikey string) (*Client, error) {
	// Parse WebSocket URL
	addr := os.Getenv("WEBSOCKET_URL")
	if addr == "" {
		addr = "wss://sensory-6d32.onrender.com/ws/logs"
		//fmt.Println("WEBSOCKET_URL not set, using default:", addr)
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL '%s': %w", addr, err)
	}

	// Connect via WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}

	client := &Client{
		conn:   conn,
		ApiKey: apikey,
		noOp:   false,
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

	// Send via WebSocket
	if c.conn != nil {
		err := c.conn.WriteJSON(entry)
		if err != nil {
			fmt.Println("WebSocket write failed:", err)
			return err
		}
	}
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

/*package logClient

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

func New(apikey, addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:   conn,
		ApiKey: apikey,
		noOp:   false,
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
}*/
