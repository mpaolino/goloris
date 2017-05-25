package main

import (
	"math/rand"
	"time"
)

// Browser constants.
const (
	_       = iota
	Chrome  = iota
	Firefox = iota
	IE      = iota
	Opera   = iota
	Safari  = iota
)

// Spec defines the user agent to be generated.
type Spec struct {
	Browser int
}

// Mock list of browser strings to randomly "generate" one from
var mockUAs = map[int][]string{
	Chrome: []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/26.0.1234.56 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/25.0.1234.56 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/24.0.1234.56 Safari/537.36",
	},

	Firefox: []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64; x64; rv:22.0) Gecko/20100101 Firefox/22.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; x64; rv:21.0) Gecko/20100101 Firefox/21.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; x64; rv:20.0) Gecko/20100101 Firefox/20.0",
	},

	IE: []string{
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0)",
	},

	Opera: []string{
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.2.15 Version/10.00",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.7.62 Version/11.01",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.12.388 Version/12.14",
	},

	Safari: []string{
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.20.3 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.3",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.19.2 (KHTML, like Gecko) Version/5.0.3 Safari/533.19.2",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.18.1 (KHTML, like Gecko) Version/5.0.2 Safari/533.18.1",
	},
}

// GenerateUA create a user agent using the given specification.
func GenerateUA(s Spec) string {
	if 0 == s.Browser {
		rand.Seed(time.Now().UnixNano())
		s.Browser = 1 + rand.Intn(5)
	}

	rand.Seed(time.Now().UnixNano())

	return mockUAs[s.Browser][rand.Intn(len(mockUAs[s.Browser]))]
}

// GenerateUA create a user agent using the given specification.
func GenerateRandomUA() string {
	rand.Seed(time.Now().UnixNano())
	var browser int = 1 + rand.Intn(5)
	return mockUAs[browser][rand.Intn(len(mockUAs[browser]))]
}
