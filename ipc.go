package runit_svm

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/isbm/runit-svm/rsvc"
)

type IPCServer struct {
	listener   *net.UnixListener
	services   *rsvc.SvmServices
	socketPath string
	socketDir  string
}

func NewIPCServer(services *rsvc.SvmServices) *IPCServer {
	ipc := new(IPCServer)
	ipc.socketDir = "/tmp"
	ipc.socketPath = path.Join(ipc.socketDir, "inid.socket")
	ipc.services = services
	ipc.connect()

	return ipc
}

func (ipc *IPCServer) connect() {
	_, err := os.Stat(ipc.socketPath)
	if err == nil {
		os.Remove(ipc.socketPath)
	}

	ipc.listener, err = net.ListenUnix("unix", &net.UnixAddr{Name: ipc.socketPath, Net: "unix"})
	if err != nil {
		log.Printf("Unable to start IPC: %s. Switching to deaf mode.", err.Error())
	}
}

func (ipc *IPCServer) send(conn *net.UnixConn) {
	var buf [1024]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Printf("IO error: %s", err.Error())
		conn.Close()
	} else {
		if _, err := conn.Write([]byte(ipc.dispatch(strings.TrimSpace(string(buf[:n]))))); err != nil {
			log.Printf("Error sending IPC response: %s", err.Error())
		}
		conn.Close()
	}

}

func (ipc *IPCServer) dispatch(command string) string {
	cmd := strings.Split(command, " ")

	switch cmd[0] {
	case "list":
		return ipc.list()
	case "status":
		return ipc.status(cmd[1])
	case "stop":
		return ipc.stop(cmd[1])
	case "start":
		return ipc.start(cmd[1])
	case "restart":
		// XXX: needs errors, so no need to repeat things twice
		return ipc.stop(cmd[1]) + "\n" + ipc.start(cmd[1])
	default:
		return "Not implemented yet"
	}
}

func (ipc *IPCServer) status(name string) string {
	service, err := ipc.services.GetServiceByName(name)
	if err != nil {
		return fmt.Sprintf("Service %s error: %s", name, err.Error())
	}

	var status string
	if len(service.GetProcesses()) > 0 {
		status = "running"
	} else {
		status = "inactive"
	}

	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%s - %s\n", service.GetServiceConfiguration().GetName(), service.GetServiceConfiguration().Info))
	buff.WriteString(fmt.Sprintf("Status: %s\n", status))
	for _, p := range service.GetProcesses() {
		buff.WriteString(fmt.Sprintf("\\_ %d\n", p.Getpid()))
	}

	return buff.String()
}

func (ipc *IPCServer) start(name string) string {
	service, err := ipc.services.GetServiceByName(name)
	if err != nil {
		return fmt.Sprintf("Service %s error: %s", name, err.Error())
	}
	if err := service.Start(); err != nil {
		return fmt.Sprintf("Error starting service %s: %s", service.GetServiceConfiguration().GetName(), err.Error())
	}

	return fmt.Sprintf("Service %s has been started", service.GetServiceConfiguration().GetName())
}

func (ipc *IPCServer) stop(name string) string {
	service, err := ipc.services.GetServiceByName(name)
	if err != nil {
		return fmt.Sprintf("Service %s error: %s", name, err.Error())
	}
	if err := service.Stop(); err != nil {
		return fmt.Sprintf("Error stopping service %s: %s", service.GetServiceConfiguration().GetName(), err.Error())
	}

	return fmt.Sprintf("Service %s has been stopped", service.GetServiceConfiguration().GetName())
}

func (ipc *IPCServer) list() string {
	var buff bytes.Buffer
	for idx, runlevel := range ipc.services.GetRunlevels() {
		buff.WriteString(fmt.Sprintf("Stage %d\n", idx+1))

		for _, s := range runlevel.GetServices() {
			buff.WriteString(fmt.Sprintf("  \\_ %s - %s\n", s.GetServiceConfiguration().GetName(), s.GetServiceConfiguration().Info))
		}
	}
	return buff.String()
}

func (ipc *IPCServer) ServeForever() {
	log.Printf("--- Final stage")
	if ipc.listener != nil {
		defer os.Remove(ipc.socketPath)
		for {
			conn, err := ipc.listener.AcceptUnix()
			if err != nil {
				log.Printf("Listener error: %s", err.Error())
				conn.Close()
				break
			} else {
				go ipc.send(conn)
			}
		}
	}

	log.Printf("IPC failed. Running in deaf mode (no sysctl)")
	for {
		select {}
	}

}
