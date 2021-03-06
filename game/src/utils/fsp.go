package utils

import (
	"bufio"
	"io"
	"log"
	"net"
)

const (
	DefaultPolicy = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>\x00"
)

type FspServer struct {
	listener      *net.TCPListener
	policyContent string
}

func main() {
	fsp := NewFspServer()
	fsp.Startup()
}

func NewFspServer() *FspServer {
	fsp := &FspServer{}
	fsp.policyContent = DefaultPolicy
	return fsp
}

func (fs *FspServer) Startup() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalln(err)
	}
	fs.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening on %s", addr.String())

	for {
		conn, err := fs.listener.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}
		go fs.read(conn)
	}
	fs.listener.Close()
}


func (fs *FspServer) read(conn *net.TCPConn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		_, err := r.ReadString('\x00')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		w := bufio.NewWriter(conn)
		w.WriteString(fs.policyContent)
		w.Flush()
		log.Printf("Sent policy file to %s", conn.RemoteAddr().String())
	}
}

