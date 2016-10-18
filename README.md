# memcachedstress
A simple memcached connection stress test written in Go

Usage:
```
Usage of C:\workspace\memcache-stress\memcachedstress.exe:
  -conn_num int
        concurrent connections (default 10)
  -duration int
        time to run test in seconds (default 10)
  -host string
        hostname or ip (default "localhost")
  -port string
        memcached port (default "11211")
  -timeout int
        client timeout in millisecond's (default 100)

example
-------
mecachedstress -host=localhost -port=11211 -conn_num-500 -duration=30

```
