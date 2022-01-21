package vrpc

import (
	"bytes"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/snowmerak/vrpc/frame"
)

type Client struct {
	sync.Mutex
	logger *log.Logger
	buf    *bytes.Buffer
	conn   net.Conn
	addr   string
}

func NewClient(addr string, logger *log.Logger) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		logger: logger,
		conn:   conn,
		buf:    bytes.NewBuffer(nil),
		addr:   addr,
	}
	client.logger.Printf("%s: connected\n", conn.RemoteAddr())
	return client, nil
}

func (c *Client) Request(service uint32, method uint32, data []byte) ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	if c.conn == nil {
		return nil, errors.New("client is not connected")
	}
	data = frame.New_Frame(service, method, 0, uint32(len(data)), data)
	n, err := c.conn.Write(data)
	if err != nil {
		c.logger.Printf("%s: write error: %s\n", c.conn.RemoteAddr(), err)
		return nil, err
	}
	if n != len(data) {
		c.logger.Printf("%s: write error: %s\n", c.conn.RemoteAddr(), "short write")
		return nil, errors.New("failed to write all data")
	}
	buf := [1024]byte{}
	header := [16]byte{}
	n, err = c.conn.Read(header[:])
	if err != nil {
		c.logger.Printf("%s: read error: %s\n", c.conn.RemoteAddr(), err)
		return nil, err
	}
	if n != len(header) {
		c.logger.Printf("%s: read error: %s\n", c.conn.RemoteAddr(), "short read")
		return nil, errors.New("failed to read all data")
	}
	bodySize := uint64(frame.Frame(header[:]).BodySize())
	curSize := uint64(0)
	c.buf.Reset()
	for {
		n, err := c.conn.Read(buf[:])
		if err != nil {
			c.logger.Printf("%s: error occurred: %s\n", c.conn.RemoteAddr(), err)
			return nil, err
		}
		c.buf.Write(buf[:n])
		curSize += uint64(n)
		if curSize >= bodySize {
			c.logger.Printf("%s: read %d bytes\n", c.conn.RemoteAddr(), curSize)
			return c.buf.Bytes(), nil
		}
	}
}

func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.logger.Printf("%s: disconnected\n", c.conn.RemoteAddr())
		c.conn = nil
	}
}

func (c *Client) Reconnect() error {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		c.logger.Printf("%s: reconnect error: %s\n", c.addr, err)
		return err
	}
	c.conn = conn
	c.logger.Printf("%s: reconnected\n", c.conn.RemoteAddr())
	return nil
}
