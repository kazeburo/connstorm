package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/valyala/tcplisten"
)

type commonSetting struct {
	FromPort        uint `long:"from" description:"lower port number to listen or connect"`
	ToPort          uint `long:"to" description:"upper port number to listen or connect"`
	newConnections  uint64
	prevConnections uint64
}

type srvOpts struct {
	commonSetting
	BindAddr string        `long:"listen" short:"l" default:"0.0.0.0" description:"address for listen"`
	Linger   int           `long:"linger" default:"0" description:"lingering timeout"`
	Delay    time.Duration `long:"delay" default:"0.1s" description:"delay time before close socket"`
}

type cliOpts struct {
	commonSetting
	Addrs       []string      `long:"host" short:"H" required:"true" description:"hostname[s] to connect"`
	MaxWorkers  uint          `long:"max-workers" default:"100" description:"max number of worker to connect to server"`
	ReadTimeout time.Duration `long:"read-timeout" default:"30s" description:"read timeout"`
}

var globalPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 1024)
		return &b
	},
}

func (opts *srvOpts) handleConnection(conn net.Conn) {
	defer conn.Close()
	atomic.AddUint64(&opts.newConnections, 1)
	conn.(*net.TCPConn).SetLinger(opts.Linger)
	time.Sleep(opts.Delay)
}

func (opts *srvOpts) handleListener(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Temporary() {
					log.Printf("AcceptTCP: %v", err)
					continue
				}
			}
			log.Fatal("Accept: %v", err)
			return
		}
		go opts.handleConnection(conn)
	}
}

func (opts *srvOpts) Counter() {
	for {
		time.Sleep(10 * time.Second)
		new := opts.newConnections
		log.Printf("newConnection: %f connections/sec", (float64(new)-float64(opts.prevConnections))/float64(10))
		opts.prevConnections = new
	}
}

func (opts *srvOpts) Execute(args []string) error {
	if opts.FromPort > opts.ToPort {
		return fmt.Errorf("--to is lower than --from")
	}
	opts.newConnections = 0
	opts.prevConnections = 0

	lc := tcplisten.Config{
		DeferAccept: false,
		FastOpen:    false,
		ReusePort:   true,
	}

	for port := opts.FromPort; port <= opts.ToPort; port++ {
		l, err := lc.NewListener("tcp4", net.JoinHostPort(opts.BindAddr, fmt.Sprintf("%d", port)))

		if err != nil {
			return err
		}
		go opts.handleListener(l)
	}
	opts.Counter()
	return nil
}

var dialer = net.Dialer{
	Timeout:       10 * time.Second,
	FallbackDelay: -1 * time.Second,
	KeepAlive:     -1 * time.Second,
}

func (opts *cliOpts) cliWorker(addr string) {
	conn, err := dialer.Dial("tcp4", addr)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer conn.Close()
	atomic.AddUint64(&opts.newConnections, 1)
	conn.SetReadDeadline(time.Now().Add(opts.ReadTimeout))
	buf := globalPool.Get().(*[]byte)
	defer func() {
		globalPool.Put(buf)
	}()
	conn.Read(*buf)
}

func (opts *cliOpts) Counter() {
	for {
		time.Sleep(10 * time.Second)
		new := opts.newConnections
		log.Printf("newConnection: %f connections/sec", (float64(new)-float64(opts.prevConnections))/float64(10))
		opts.prevConnections = new
	}
}

func (opts *cliOpts) Execute(args []string) error {
	if opts.FromPort > opts.ToPort {
		return fmt.Errorf("--to is lower than --from")
	}
	opts.newConnections = 0
	opts.prevConnections = 0
	addrs := make([]string, 0)

	for port := opts.FromPort; port <= opts.ToPort; port++ {
		for _, addr := range opts.Addrs {
			tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(addr, fmt.Sprintf("%d", port)))
			if err != nil {
				return err
			}
			addrs = append(addrs, tcpAddr.String())
		}
	}
	ch := make(chan string, opts.MaxWorkers*2)
	for w := uint(0); w < opts.MaxWorkers; w++ {
		go func() {
			for {
				a := <-ch
				opts.cliWorker(a)
			}
		}()
	}
	go func() {
		for {
			for _, a := range addrs {
				ch <- a
			}
		}
	}()
	opts.Counter()
	return nil
}

type mainOpts struct {
	ServerCmd srvOpts `command:"server"`
	ClientCmd cliOpts `command:"client"`
}

func main() {
	opts := mainOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	_, err := psr.Parse()
	if err != nil {
		os.Exit(1)
	}
}
