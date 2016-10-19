# memcachedstress
A simple memcached connection stress test written in Go

```
Usage of ./memcachedstress:
  -conn_num int
    	concurrent connections (default 10)
  -duration int
    	time to run test (default 10)
  -host string
    	hostname or ip (default "localhost")
  -log_level int
    	how much logging to display levels 1 - 6 with 1 being the most verbose and 6 the least (default 2)
  -port string
    	memcached port (default "11211")
  -ramp int
    	time between new connection spawn in millisecond's (default 50)
  -sleep int
    	time in seconds between new set and get in spawned connections (default 2)
  -timeout int
    	client timeout in millisecond's (default 100)

example
-------
mecachedstress -host=localhost -port=11211 -conn_num-500 -duration=30

```
