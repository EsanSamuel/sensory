package models

import "net"

type Client struct {
	conn      net.Conn
	Project   string
	Service   string
	ApiKey    string
	ProjectId string
	UserId    string
	noOp      bool // flag to mark dummy client
}
