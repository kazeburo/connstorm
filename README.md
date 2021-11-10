# connstorm

make TCP connection storm between server and client for benchmarking network stuff.

## usage

```
Usage:
  connstorm [OPTIONS] <client | server>

Help Options:
  -h, --help  Show this help message

Available commands:
  client
  server
```


client help

```
Usage:
  connstorm [OPTIONS] client [client-OPTIONS]

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
  connstorm [OPTIONS] server [server-OPTIONS]

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
$ /usr/local/bin/connstorm server --from 8500 --to 8800 --linger 1 --delay 0.1s
```

run client 


```
$ GOGC=500 /usr/local/bin/connstorm client -H server1 -H server2 -H server3 --from 8500 --to 8800 --max-workers 10000
2021/11/09 13:12:04 newConnection: 103270.200000 connections/sec
2021/11/09 13:12:14 newConnection: 105419.600000 connections/sec
2021/11/09 13:12:24 newConnection: 108597.700000 connections/sec
```

