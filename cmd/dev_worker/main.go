package main

import (
	"fmt"
	"os"
	"os/signal"
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

type Context struct {
	customerID int64
}

func main() {
	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "my_app_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 10, "my_app_namespace", redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCustomer)

	// Map the name of jobs to handler functions
	//pool.Job("send_email", (*Context).SendEmail)

	// Customize options:
	pool.JobWithOptions("send_email", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).SendEmail)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) FindCustomer(job *work.Job, next work.NextMiddlewareFunc) error {
	// If there's a customer_id param, set it in the context for future middleware and handlers to use.
	if _, ok := job.Args["customer_id"]; ok {
		c.customerID = job.ArgInt64("customer_id")
		if err := job.ArgError(); err != nil {
			return err
		}
	}

	return next()
}

func (c *Context) SendEmail(job *work.Job) error {
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	// Go ahead and send the email...
	// sendEmailTo(addr, subject)

	fmt.Printf("add : %s, subject : %s\n", addr, subject)

	return fmt.Errorf("lohhhhh")
}

func (c *Context) Export(job *work.Job) error {
	return nil
}
