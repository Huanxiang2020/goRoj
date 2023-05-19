package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	PoolUnInvaildSizeError = errors.New("pool size is unvaild")
	PoolIsClosedError      = errors.New("pool had closed")
)

// 连接池定义
type PoolMySQL struct {
	sync.Mutex                // 保证连接池线程安全
	Size       int            // 连接池连接数量
	ConnChan   chan io.Closer // 存储连接的管道
	IsClose    bool
	ctx        context.Context
}

func NewMySQLPool(size int) (*PoolMySQL, error) {
	if size <= 0 {
		return nil, PoolUnInvaildSizeError
	}
	return &PoolMySQL{
		ConnChan: make(chan io.Closer, size),
		ctx:      context.Background(),
	}, nil
}

// 获取连接
func (pool *PoolMySQL) GetConnFromPool() (io.Closer, error) {
	if pool.IsClose == true {
		return nil, PoolIsClosedError
	}
	select {
	// 从管道中消费
	case conn, ok := <-pool.ConnChan:
		if !ok {
			return nil, PoolIsClosedError
		}
		fmt.Println("获取到连接：", conn)
		return conn, nil
	default:
		// 连接池中没有连接，新建连接
		return pool.getNewConn(pool.ctx)
	}
}

// 构造新连接
func (pool *PoolMySQL) getNewConn(ctx context.Context) (io.Closer, error) {
	db, err := sql.Open("mysql", "root:chenyang@tcp(127.0.0.1:3306)/customers?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal("数据库连接失败", err)
		return nil, err
	}
	conn, _ := db.Conn(ctx)
	select {
	case pool.ConnChan <- conn:
		fmt.Println("连接放入连接池")
	default:
		fmt.Println("连接池满了，连接丢弃")
		conn.Close()
	}
	return conn, nil
}

// 释放连接
func (pool *PoolMySQL) ReleaseConn(conn io.Closer) error {
	pool.Lock()
	defer pool.Unlock()

	if pool.IsClose == true {
		return PoolIsClosedError
	}

	select {
	case pool.ConnChan <- conn:
		fmt.Println("连接已放回", conn)
	default:
		fmt.Println("连接池满了，连接丢弃")
		conn.Close()
	}
	return nil
}

// 关闭连接池
func (pool *PoolMySQL) ClosePool() error {
	pool.Lock()
	defer pool.Unlock()

	if pool.IsClose == true {
		return PoolIsClosedError
	}

	pool.IsClose = true
	close(pool.ConnChan)

	for conn := range pool.ConnChan {
		conn.Close()
	}
	return nil
}

func TestMySQLConn() {
	wg := sync.WaitGroup{}
	pool, err := NewMySQLPool(16)
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(20)
	fmt.Println("开启20个协程获取连接")
	for i := 0; i < 20; i++ {
		go TestReleaseAndGetConn(pool, wg)
	}
	wg.Wait()
	fmt.Println("main end")
}

func TestReleaseAndGetConn(pool *PoolMySQL, wg sync.WaitGroup) {
	s := rand.Int63n(2)
	time.Sleep(time.Duration(s) * time.Second)
	conn, err := pool.GetConnFromPool()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接池连接数：", len(pool.ConnChan))
	time.Sleep(time.Duration(s) * time.Second)
	pool.ReleaseConn(conn)
	wg.Done()
}
