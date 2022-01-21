package vrpc

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"net"
	"sync"

	"github.com/snowmerak/vrpc/frame"
)

type Client struct {
	lock   sync.Mutex
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
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn == nil {
		return nil, errors.New("client is not connected")
	}
	size := uint32(len(data))
	seq := rand.Uint32()
	data = frame.New_Frame(service, method, seq, size+8, data)
	n, err := c.conn.Write(data)
	if err != nil {
		c.logger.Printf("%s: write error: %s\n", c.conn.RemoteAddr(), err)
		return nil, err
	}
	if n != len(data) {
		c.logger.Printf("%s: write error: %s\n", c.conn.RemoteAddr(), "short write")
		return nil, errors.New("failed to write all data")
	}
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
	buf := [1024]byte{}
	c.buf.Reset()
	c.buf.Write(header[:])
	for {
		n, err := c.conn.Read(buf[:])
		if err != nil {
			c.logger.Printf("%s: error occurred: %s\n", c.conn.RemoteAddr(), err)
			return nil, err
		}
		c.buf.Write(buf[:n])
		curSize += uint64(n)
		if curSize >= bodySize {
			if curSize > bodySize {
				c.logger.Printf("%s: body size mismatched %d > %d\n", c.conn.RemoteAddr(), curSize, bodySize)
				return nil, errors.New("body size is larger than expected")
			}
			data := frame.Frame(c.buf.Bytes())
			if !data.Vstruct_Validate() {
				c.logger.Printf("%s: invalid frame arrived\n", c.conn.RemoteAddr())
				return nil, errors.New("invalid frame arrived")
			}
			if data.Sequence() != seq+1 {
				c.logger.Printf("%s: sequence mismatched %d != %d\n", c.conn.RemoteAddr(), data.Sequence(), seq+1)
				return nil, errors.New("sequence mismatched")
			}
			c.logger.Printf("%s: read %d bytes\n", c.conn.RemoteAddr(), curSize)
			return data.Body(), nil
		}
	}
}

func (c *Client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.logger.Printf("%s: disconnected\n", c.conn.RemoteAddr())
		c.conn = nil
	}
}

func (c *Client) Reconnect() error {
	c.lock.Lock()
	defer c.lock.Unlock()
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
