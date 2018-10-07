package main

import (
	"bufio"
	//  "fmt"
	"net"
	"os"
	"strings"
	"time"
	//  "strconv"

	"github.com/sirupsen/logrus"
	"github.com/skillcoder/homer/shutdown"
)

var log = logrus.New()

func sendLineToClient(c net.Conn, str string) {
	c.Write([]byte(str + "\n"))
	addr := c.RemoteAddr().String()
	log.Debugf(">%s %s", addr, str)
}

func handleConnection(c net.Conn) {
	addr := c.RemoteAddr().String()
	log.Infof("Connected %s", addr)
	// Close connection when this function ends
	defer func() {
		log.Info("Closing connection", addr)
		c.Close()
	}()
	bufReader := bufio.NewReader(c)
	timeoutDuration := 10 * time.Second
	for {
		// Set a deadline for reading. Read operation will fail if no data
		// is received after deadline.
		c.SetReadDeadline(time.Now().Add(timeoutDuration))

		netData, err := bufReader.ReadString('\n')
		if err != nil {
			log.Error(err)
			return
		}

		cmd := strings.TrimSpace(netData)
		log.Debugf("<%s %s", addr, cmd)
		if cmd == "exit" {
			log.Info("CMD EXIT", addr)
			sendLineToClient(c, "Bye!")
			break
		}

		//result := strconv.Itoa(time) + "\n"
		c.Write([]byte(netData))
		log.Debugf(">%s %s", addr, strings.TrimSpace(string(netData)))
	}
}

func tcp_listening(l net.Listener) {
	for {
		c, err := l.Accept()
		//log.Debug("Accept")
		if err != nil {
			log.Error(err)
			return
		}
		go handleConnection(c)
	}
}

func init() {
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	tcp_listen := os.Getenv("TCPING_SERVICE_LISTEN")
	if len(tcp_listen) == 0 {
		log.Fatal("Required env parameter TCPING_SERVICE_LISTEN [ip:port] is not set")
	}

	l, err := net.Listen("tcp4", tcp_listen)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		l.Close()
		log.Info("Listener closed")
	}()
	log.Infof("Listen on %s", tcp_listen)

	go tcp_listening(l)

	logger := log.WithField("event", "shutdown")
	sdHandler := shutdown.NewHandler(logger)
	sdHandler.RegisterShutdown(sd)
}

// sd does graceful dhutdown of the service
func sd() (string, error) {
	// if service has to finish some tasks before shutting down, these tasks must be finished her
	// TODO(developer): wait for all gorutined ends
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
	return "Ok", nil
}
