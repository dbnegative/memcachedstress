package main

import (
	"flag"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bradfitz/gomemcache/memcache"
)

var host = flag.String("host", "localhost", "hostname or ip")
var port = flag.String("port", "11211", "memcached port")
var num = flag.Int("conn_num", 10, "concurrent connections")
var testduration = flag.Int("duration", 10, "time to run test")
var timeout = flag.Int("timeout", 100, "client timeout in millisecond's")
var succesful int
var failure int

func main() {

	flag.Parse()
	log.SetLevel(log.InfoLevel)

	for count := *num; count > 0; count-- {
		s := strconv.Itoa(count)

		log.WithFields(log.Fields{
			"Routine Number": count,
			"Key":            "foo" + s,
			"Value":          s,
		}).Info("Spinning off goroutine")

		go connectMemcache(*host, *port, "foo"+s, s)

		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(time.Duration(*testduration) * time.Second)

	log.Printf("Waiting for %vs", *testduration)

	log.WithFields(log.Fields{
		"Conn Num":        *num,
		"Duration":        *testduration,
		"Successful Gets": succesful,
		"Failed Gets":     failure,
	}).Info("Finished Test")
}

func connectMemcache(host string, port string, key string, value string) {
	log.WithFields(log.Fields{
		"Host":    host,
		"Port":    port,
		"Key":     key,
		"Value":   value,
		"Timeout": *timeout,
	}).Info("Created new memcache connection")

	mc := memcache.New(host + ":" + port)
	mc.Timeout = time.Duration(*timeout) * time.Millisecond

	for i := 0; ; i++ {
		mc.Set(&memcache.Item{Key: key, Value: []byte(value)})
		_, err := mc.Get(key)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Could not get key")
			failure++
			err = nil //reset error
		} else {
			succesful++
		}

		time.Sleep(1 * time.Second)
	}
}
