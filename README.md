# webdav
file server written in go; serve files over http/https using webdav.

```
rxlx ~ $ webdav -h
Usage of webdav:
  -anon
    	anonymous connections allowed (user auth disabled)
  -both
    	run an http server and https server
  -cert string
    	path to your cert (default "cert.pem")
  -dir string
    	Directory to serve from. Default is CWD (default "./")
  -insecure
    	disable TLS
  -key string
    	path to your key (default "key.pem")
  -log string
    	syslog server or /path/to/file to log to (default "webdav.log")
  -monitor
    	enable metric logging; memory, heap, numGC, etc
  -p int
    	http port (plain) (default 6200)
  -poll int
    	how often to poll runtime stats (default 30)
  -ps int
    	https port (tls) (default 6201)
  -quiet
    	only log errors
  -uniq string
    	if using syslog, a unique process name for easier debugging (default "__DAV__")
  -v	show version number
  
 **NOTE** you'll need a cert.pem and key.pem for tls to work

# example runs in TLS only mode, allows anonymous connections, logs to a remote server, but only errors
# the uniq flag is just a way to grep through your remote logs easier. finally , run on port 6201 and monitor the
# runtime stats every 30 minutes
# EXAMPLE:
 webdav -dir ./share/ -anon -log udp@192.168.86.42:514 -quiet -uniq "MySpecialProcess01" \
        -cert ./cert.pem -key ./key.pem -ps 6201 -monitor -poll 1800
```

<h2>install / config</h2>
Basic Auth Config:
unless modified / recompiled, webdav will look for the environment variables listed below. They must be present server side for basic auth to work. The actual values can be whatever, just keep DUSR and DAT unless you know what you're doing.
<hr>

```bash
export DUSR="cowpower"                                # webdav user
export DAT="a961431417E^Cab5d9fDe53752ec81937dc944*5" # webdav access token
```

<hr>

*TO INSTALL*

if you're compiling from source, you'll need to install go (dont worry, it's easy) -> https://go.dev/doc/install. otherwise just download the binary and add it to your path.
<br>

```bash
mkdir $HOME/bin/
git clone https://github.com/rexlx/webdav.git
cd webdav/
go build webdav.go
mv webdav $HOME/bin/   # or whatever path location you'd like
# if you have no use for the source code, optionally:
cd ../;rm -rf ./webdav/;cd $HOME;ls $HOME/bin/

# finally, verify installed and on the PATH with:
webdav -h
```


<hr>
<h3>future Updates in the works</h3>
Support for json configuration files to reduce CLI args. More robust user authentication. add logging levels. 
