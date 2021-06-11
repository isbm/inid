package sysctl

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

type SvmCtl struct {
	inidSocket       string
	inidClientSocket string
}

func NewSvmCtl() *SvmCtl {
	ctl := new(SvmCtl)
	ctl.inidSocket = "/tmp/inid-cli.socket"
	ctl.inidSocket = "/tmp/inid.socket"
	return ctl
}

func (ctl *SvmCtl) Help() {
	exn := path.Base(os.Args[0])
	fmt.Printf("Usage:\n\t%s list\n\t%s [start|stop|restart|status] [SERVICE]\n", exn, exn)
}

func (ctl *SvmCtl) read(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:n]))
	}
}

func (ctl *SvmCtl) send(msg string) {
	_, err := os.Stat(ctl.inidClientSocket)
	if err == nil {
		os.Remove(ctl.inidClientSocket)
	}

	proto := "unix"
	laddr := net.UnixAddr{Name: ctl.inidClientSocket, Net: proto}
	conn, err := net.DialUnix(proto, &laddr, &net.UnixAddr{Name: ctl.inidSocket, Net: proto})
	if err != nil {
		panic(err)
	}

	defer os.Remove(ctl.inidClientSocket)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	ctl.read(conn)
	conn.Close()
}

func (ctl *SvmCtl) Run() error {
	if len(os.Args) < 2 || len(os.Args) == 2 && os.Args[1] != "list" {
		ctl.Help()
		os.Exit(0)
	}
	switch os.Args[1] {
	case "list":
		ctl.send("list")
	case "start", "stop", "restart", "status":
		ctl.send(strings.Join(os.Args[1:], " "))
	default:
		fmt.Println("Unknown option:", os.Args[1])
	}
	return nil
}
