package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
}
