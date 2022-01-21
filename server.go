package vrpc

import (
	"bytes"
	"log"
	"net"

	"github.com/snowmerak/vrpc/frame"
)

type Server struct {
	logger  *log.Logger
	conns   map[net.Conn]bool
	methods map[uint32]map[uint32]func([]byte) []byte
}

func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger:  logger,
		conns:   make(map[net.Conn]bool),
		methods: make(map[uint32]map[uint32]func([]byte) []byte),
	}
}

func (s *Server) Register(service uint32, method uint32, f func([]byte) []byte) {
	if _, ok := s.methods[service]; !ok {
		s.methods[service] = make(map[uint32]func([]byte) []byte)
	}
	s.methods[service][method] = f
}

func (s *Server) Unregister(service uint32, method uint32) {
	delete(s.methods[service], method)
}

func (s *Server) handler(conn net.Conn) {
	defer func() {
		conn.Close()
		delete(s.conns, conn)
		s.logger.Printf("%s disconnected\n", conn.RemoteAddr())
	}()
	buf := [1024]byte{}
	data := bytes.NewBuffer(nil)

	s.conns[conn] = true
	s.logger.Printf("%s: connected\n", conn.RemoteAddr())
	header := [16]byte{}
	n, err := conn.Read(header[:])
	if err != nil {
		s.logger.Printf("%s: cannot read header: %s\n", conn.RemoteAddr(), err)
		return
	}
	if n != len(header) {
		s.logger.Printf("%s: cannot read header: %s\n", conn.RemoteAddr(), "short read")
		return
	}
	bodySize := uint64(frame.Frame(header[:]).BodySize())
	curSize := uint64(0)
	data.Write(header[:])
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			s.logger.Printf("%s: error occurred: %s\n", conn.RemoteAddr(), err)
			data.Reset()
			continue
		}
		s.logger.Printf("%s: read %d bytes\n", conn.RemoteAddr(), n)
		data.Write(buf[:n])
		curSize += uint64(n)
		if curSize >= bodySize {
			frm := frame.Frame(data.Bytes())
			if !frm.Vstruct_Validate() {
				s.logger.Printf("%s: Invalid vstruct frame\n", conn.RemoteAddr())
				return
			}
			replyBody := s.methods[frm.Service()][frm.Method()](frm.Body())
			result := frame.New_Frame(frm.Service(), frm.Method(), frm.Sequence()+1, uint32(len(replyBody)), replyBody)
			n, err := conn.Write(result)
			if err != nil {
				s.logger.Printf("%s: error writing to connection: %s\n", conn.RemoteAddr(), err)
				return
			}
			s.logger.Printf("%s: wrote %d bytes\n", conn.RemoteAddr(), n)
		}
	}
}

func (s *Server) Serve(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Println(err)
			continue
		}

		go s.handler(conn)
	}
}

func (s *Server) Shutdown() {
	for conn := range s.conns {
		conn.Close()
	}
}
