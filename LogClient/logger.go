package logClient

import (
	"fmt"
	"os"
)

func (c *Client) INFO(msg any) {
	c.Send("INFO", fmt.Sprint(msg))
}

func (c *Client) ERROR(msg any) {
	c.Send("ERROR", fmt.Sprint(msg))
}

func (c *Client) WARN(msg any) {
	c.Send("WARN", fmt.Sprint(msg))
}

func (c *Client) DEBUG(msg any) {
	c.Send("DEBUG", fmt.Sprint(msg))
}

func (c *Client) FATAL(msg any) {
	c.Send("FATAL", fmt.Sprint(msg))
	os.Exit(1)
}
