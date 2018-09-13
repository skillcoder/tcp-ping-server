package main

import (
  "bufio"
//  "fmt"
  "net"
  "os"
//  "strconv"
  "strings"
//  "time"
  "github.com/sirupsen/logrus"
  "github.com/skillcoder/homer/shutdown"
)

var log = logrus.New()

func handleConnection(c net.Conn) {
  log.Infof("Connected %s", c.RemoteAddr().String())
  addr := c.RemoteAddr().String()
  for {
    netData, err := bufio.NewReader(c).ReadString('\n')
    if err != nil {
      log.Error(err)
      return
    }

    cmd := strings.TrimSpace(string(netData))
	log.Debugf("<%s %s", addr, cmd)
    if cmd == "exit" {
      break
    }

    //result := strconv.Itoa(time) + "\n"
    c.Write([]byte(netData))
	log.Debugf(">%s %s", addr, strings.TrimSpace(string(netData)))
  }

  log.Infof("Disconnected %s", c.RemoteAddr().String())
  c.Close()
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
  //log.SetLevel(logrus.DebugLevel)
}

func main() {
  tcp_listen := os.Getenv("TCPPINGSERVER_SERVICE_LISTEN")
  if len(tcp_listen) == 0 {
    log.Fatal("Required env parameter TCPPINGSERVER_SERVICE_LISTEN [ip:port] is not set")
  }

  l, err := net.Listen("tcp4", tcp_listen)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer l.Close()
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

