# session-overflow

make connection storm between server and client for benchmarking network stuff.

## usage

```
Usage:
  session-overflow [OPTIONS] <client | server>

Help Options:
  -h, --help  Show this help message

Available commands:
  client
  server
```


client help

```
Usage:
  session-overflow [OPTIONS] client [client-OPTIONS]

Help Options:
  -h, --help              Show this help message

[client command options]
          --from=         lower port number to listen or connect
          --to=           upper port number to listen or connect
      -H, --host=         hostname[s] to connect
          --max-workers=  max number of worker to connect to server (default: 100)
          --read-timeout= read timeout (default: 30s)
```

server help

```
Usage:
  session-overflow [OPTIONS] server [server-OPTIONS]

Help Options:
  -h, --help        Show this help message

[server command options]
          --from=   lower port number to listen or connect
          --to=     upper port number to listen or connect
      -l, --listen= address for listen (default: 0.0.0.0)
          --linger= lingering timeout (default: 0)
          --delay=  delay time before close socket (default: 0.1s)
```


## example

exec server

```
$ /usr/local/bin/session-overflow server --from 8500 --to 8800 --linger 1 --delay 0.1s
```

run client 


```
$ GOGC=500 /usr/local/bin/session-overflow client -H server1 -H server2 -H server3 --from 8500 --to 8800 --max-workers 10000
```

