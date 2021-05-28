# webdav
small webdav server written in go. meant for network attached storage.

webdav is a technology related to http. Got tired of SMB / NFS compatibility issues and decided to give webdav a try. This is a pretty vanilla recipe, nothing mind blowing here.

I'd like to add user support but I will need to implement some middleware for that. Another thing to note: this code creates a memory lock to prevent it from paging. sometimes on sigint it'll trigger a wait time that can last a few minutes. you can safely kill it with 9.

```
rxlx ~ $ webdav -h
Usage of webdav:
  -d string
    	Directory to serve from. Default is CWD (default "./")
  -l string
    	path/file to log to (default "./webdav.log")
  -p int
    	Port to serve on (Plain HTTP) (default 8081)
  -ps int
    	Port to serve TLS on (default 8443)
  -s	Serve HTTPS. Default false
  
  nohup /home/rxlx/bin/webdav -s -d /Bstor/ -ps 5201 -l /home/rxlx/bin/logs/$(date | tr ' ' '-').log &
  
 **NOTE** you'll need a cert.pem and key.pem for tls to work
```

<h2>more notes</h2>

For some reason (as of big sur 11.3.1) macOS fails to negotiate for cipher suites and the handshake will always fail for self-signed certs. You will need to add them to your keyring. **Additionally** on mac, I'm tracing down a an error where finder locks up trying to write to the remote location ( which is running on EL8 ) over http (not secure). macOS just doesnt like dealing with plain http very much it would appear.

The binary was compiled on RHEL8, but I believe it should be pretty portable.
