package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var host = flag.String("host", "localhost", "hostname or ip")
var port = flag.String("port", "11211", "memcached port")
var num = flag.Int("conn_num", 10, "concurrent connections")
var testduration = flag.Int("duration", 10, "time to run test")
var timeout = flag.Int("timeout", 100, "client timeout in millisecond's")

func main() {

	flag.Parse()

	for count := *num; count > 0; count-- {
		s := strconv.Itoa(count)
		log.Printf("Spinning up goroutine %s", s)

		go connectMemcache("localhost", "11211", "foo"+s, s)

		time.Sleep(50 * time.Millisecond)
	}

	log.Printf("Spinning up goroutine %s", "FINISHED")

	time.Sleep(time.Duration(*testduration) * time.Second)
}

func connectMemcache(host string, port string, key string, value string) {
	log.Printf("Process(%s) Created Connection: %s with value %s", value, host, port)

	mc := memcache.New(host + ":" + port)
	mc.Timeout = time.Duration(*timeout) * time.Millisecond

	for i := 0; ; i++ {

		mc.Set(&memcache.Item{Key: key, Value: []byte(value)})
		log.Printf("Process(%s) Searching for key: %s", value, key)
		_, err := mc.Get(key)

		if err != nil {
			log.Printf("Process(%s) ERROR: Could not find key: %v ", value, err)
		}

		log.Printf("Process(%s) Found Key: %s with value: %s ..sleep %d seconds", value, key, value, *timeout/1000)
		time.Sleep(1 * time.Second)
	}
}
