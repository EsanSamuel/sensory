package logClient

import "os"

func (c *Client) INFO(msg string) {
	c.Send("INFO", msg)
}

func (c *Client) ERROR(msg string) {
	c.Send("ERROR", msg)
}

func (c *Client) FATAL(msg string) {
	c.Send("FATAL", msg)
	os.Exit(1)
}
