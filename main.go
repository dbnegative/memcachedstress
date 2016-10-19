package main

import (
	"flag"
	"math/rand"
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
var ramp = flag.Int("ramp", 50, "time between new connection spawn in millisecond's")
var gap = flag.Int("sleep", 2, "time in seconds between new set and get in spawned connections")
var verbose = flag.Int("log_level", 2, "how much logging to display levels 1 - 6 with 1 being the most verbose and 6 the least")

var succesful int
var failure int

func main() {

	flag.Parse()

	switch *verbose {
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.ErrorLevel)
	case 5:
		log.SetLevel(log.FatalLevel)
	case 6:
		log.SetLevel(log.PanicLevel)
	}

	log.WithFields(log.Fields{
		"Connection Count": *num,
	}).Info("Spawning Connections...")

	for count := *num; count > 0; count-- {
		s := strconv.Itoa(count)
		key := generateRandKey(18)
		log.WithFields(log.Fields{
			"Routine Number": count,
			"Key":            key,
			"Value":          s,
		}).Debug("Spinning off goroutine")

		//spawn  routine
		go connectMemcache(*host, *port, key, s)

		time.Sleep(time.Duration(*ramp) * time.Millisecond)
	}

	log.WithFields(log.Fields{
		"Duration": *testduration,
	}).Info("Waiting...")

	time.Sleep(time.Duration(*testduration) * time.Second)

	log.WithFields(log.Fields{
		"Conn Num":        *num,
		"Duration":        *testduration,
		"Successful Gets": succesful,
		"Failed Gets":     failure,
		"Timeout":         *timeout,
	}).Info("Finished Test")
}

func generateRandKey(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func connectMemcache(host string, port string, key string, value string) {
	log.WithFields(log.Fields{
		"Host":    host,
		"Port":    port,
		"Timeout": *timeout,
	}).Debug("Created new memcache connection")

	mc := memcache.New(host + ":" + port)
	mc.Timeout = time.Duration(*timeout) * time.Millisecond

	for i := 0; ; i++ {
		mc.Set(&memcache.Item{Key: key, Value: []byte(value)})
		_, err := mc.Get(key)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"Key":   key,
			}).Error("Could not get key")
			failure++
			err = nil //reset error
		} else {
			succesful++
		}

		time.Sleep(time.Duration(*gap) * time.Second)
	}
}
