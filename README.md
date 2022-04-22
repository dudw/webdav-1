# webdav
small webdav server written in go. meant for network attached storage. Support for json configuration file coming soon.

```
$ webdav -h
Usage of webdav:
  -cert string
    	path to your cert (default "./cert.pem")
  -dir string
    	Directory to serve from. Default is CWD (default "./")
  -insecure
    	disable TLS
  -key string
    	path to your key (default "./key.pem")
  -log string
    	path/file to log to (default "./webdav.log")
  -monitor
    	enable metric logging; memory, heap, numGC, etc
  -p int
    	http port (plain) (default 6200)
  -poll int
    	how often to poll runtime stats (default 30)
  -ps int
    	https port (tls) (default 6201)
  
 **NOTE** you'll need a cert.pem and key.pem for tls to work

 # start the wevdav server, log to a remote syslog server, monitor stats
 # every 30 mins.
 $ webdav -log udp@ADDR:514 \
   -dir /path/here/ \
   -cert /path/cert.pem \
   -key /path/key.pem \
   -ps 5201 -monitor -poll 1800
```

<h2>post script</h2>
Running this will start both an http server as well as https. That means if you do not supply a specific port for plain http, it will default to 6021. if you have anything running on this port the server will fail to start.

<br>

<h3>more notes</h3>

For some reason (as of big sur 11.3.1) macOS fails to negotiate for cipher suites and the handshake will always fail for self-signed certs. You will need to add them to your keyring. **Additionally** on mac, I'm tracing down a an error where finder locks up trying to write to the remote location ( which is running on EL8 ) over http (not secure). macOS just doesnt like dealing with plain http very much it would appear.

The binary was compiled on RHEL8, but I believe it should be pretty portable.
