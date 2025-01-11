package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{serverAddr: address, timeout: timeout, in: in, out: out}
}

type client struct {
	serverAddr string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	conn       net.Conn
}

func (c *client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.serverAddr, c.timeout)
	return err
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	// отправить в сокет содержимое in
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *client) Receive() error {
	// прочитать доступные данные из сокета и записать их в out
	_, err := io.Copy(c.out, c.conn)
	return err
}
