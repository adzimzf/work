package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/adzimzf/work"
)

var redisPool = &redis.Pool{
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

// Make an enqueuer with a particular namespace
var enqueuer = work.NewEnqueuer("my_app_namespace", redisPool)

func main() {
	// Enqueue a job named "send_email" with the specified parameters.
	_, err := enqueuer.EnqueueUniqueInByKey("send_email", 0, work.Q{"address": "test@example.com",
		"subject":     fmt.Sprintf("jam : %s", time.Now().Format(time.RFC850)),
		"customer_id": 4}, work.Q{"address": "test@example.com",
		"subject":     fmt.Sprintf("jam : %s", time.Now().Format(time.RFC850)),
		"customer_id": 4})
	if err != nil {
		log.Fatal(err)
	}
}

type Context struct {
	customerID int64
}
