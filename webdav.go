package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/webdav"
)

const _version_ = "Rel_20220514"

var (
	httpPort  = flag.Int("p", 6200, "http port (plain)")
	httpsPort = flag.Int("ps", 6201, "https port (tls)")
	poll      = flag.Int("poll", 30, "how often to poll runtime stats")
	insecure  = flag.Bool("insecure", false, "disable TLS")
	anon      = flag.Bool("anon", false, "anonymous connections allowed (user auth disabled)")
	monitor   = flag.Bool("monitor", false, "enable metric logging; memory, heap, numGC, etc")
	both      = flag.Bool("both", false, "run an http server and https server")
	version   = flag.Bool("v", false, "show version number")
	quiet     = flag.Bool("quiet", false, "only log errors")
	cert      = flag.String("cert", "cert.pem", "path to your cert")
	key       = flag.String("key", "key.pem", "path to your key")
	dir       = flag.String("dir", "./", "Directory to serve from. Default is CWD")
	logPath   = flag.String("log", "webdav.log", "syslog server or /path/to/file to log to")
	uniq      = flag.String("uniq", "__DAV__", "if using syslog, a unique process name for easier debugging")
)

type Profile struct {
	Alloc,
	TotalAlloc,
	MemoryAlloc,
	System,
	Free,
	Objects,
	TotalPauses uint64
	NGC    uint32
	NumCPU int
}

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("\nwebdav fileserver version: %v\n", _version_)
		os.Exit(0)
	}
	// if the user supplies (what we define as a) syslog path, unpack
	// -log tcp@hostname:port | -log udp@addr:port
	if strings.Contains(*logPath, "@") {
		addr := strings.Split(*logPath, "@")
		logger, e := syslog.Dial(addr[0], addr[1],
			syslog.LOG_WARNING|syslog.LOG_DAEMON, "__DAV__") // anything else here
		check(e)
		log.SetOutput(logger)
	} else {
		// -log /this/is/sparta.log
		logger, err := os.OpenFile(*logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logger.Close()
		log.SetOutput(logger)
	}
	if *monitor {
		go monitorRuntimeProfile()
	}

	svr := &webdav.Handler{
		FileSystem: webdav.Dir(*dir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("-> %s: %s, ERROR->: %s on %v", r.Method, r.URL, err, r.RemoteAddr)
			}
			if !*quiet {
				log.Printf("-> %s: %s -> %v", r.Method, r.URL, r.RemoteAddr)
			}
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !*anon {
			uname, pwd, _ := r.BasicAuth()
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			if uname == os.Getenv("DUSR") && pwd == os.Getenv("DAT") {
				if !*quiet {
					log.Printf("recieved an authenticated connection from -> %v...starting server..", r.RemoteAddr)
				}
				w.Header().Set("Timeout", "86399")
				svr.ServeHTTP(w, r)
			} else {
				log.Printf("recieved an attempted connection from -> %v, but no credentials were provided...", r.RemoteAddr)
				w.WriteHeader(401)
				w.Write([]byte("failed to authenticate; access denied."))
			}
		} else {
			svr.ServeHTTP(w, r)
		}
	})
	if !*insecure {
		if _, err := os.Stat(*cert); err != nil {
			fmt.Printf("no cert located at: %v\n", *cert)
			os.Exit(1)
		}
		if _, er := os.Stat(*key); er != nil {
			fmt.Printf("no key located at: %v\n", *key)
			os.Exit(1)
		}
		if *both {
			go http.ListenAndServeTLS(fmt.Sprintf(":%d", *httpsPort), *cert, *key, nil)
			http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
		}
		http.ListenAndServeTLS(fmt.Sprintf(":%d", *httpsPort), *cert, *key, nil)
	}

	if *insecure {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil); err != nil {
			fmt.Println(err)
			log.Fatalf("error with webdav server (http port: %v): %v", *httpPort, err)
		}
	}
}

func monitorRuntimeProfile() {
	var p Profile
	var stats runtime.MemStats
	for {
		<-time.After(
			time.Duration(*poll) * time.Second)
		runtime.ReadMemStats(&stats)

		p.Alloc = stats.Alloc
		p.TotalAlloc = stats.TotalAlloc
		p.MemoryAlloc = stats.Mallocs
		p.Free = stats.Frees
		p.Objects = p.MemoryAlloc - p.Free

		// GC stuff
		p.NumCPU = runtime.NumCPU()
		p.NGC = stats.NumGC

		profile, _ := json.Marshal(p)
		log.Println(string(profile))
	}
}

func check(e error) {
	if e != nil {
		fmt.Printf("encountered an error!\t%v\n", e)
		os.Exit(1)
	}
}
