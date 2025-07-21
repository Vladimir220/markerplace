package db

import (
	"database/sql"
	"fmt"
	"main/tools/log"
	"os"
	"strconv"
	"sync"
)

type IConnectionPool interface {
	GetConnection() *sql.DB
	Close()
}

func CreateConnectionPool() IConnectionPool {
	tempConnectionPool := &ConnectionPool{logger: log.CreateLogger("ConnectionPool")}
	numOfConnections := os.Getenv("NUM_OF_DB_CONNECTIONS")
	if numOfConnections == "" {
		tempConnectionPool.logger.WriteError("CreateConnectionPool(): following variables is not specified in env: NUM_OF_DB_CONNECTIONS")
		return nil
	}

	i, err := strconv.ParseUint(numOfConnections, 10, 64)
	if err != nil {
		tempConnectionPool.logger.WriteError("CreateConnectionPool(): err")
		return nil
	}

	connectionPool := &ConnectionPool{
		c:      make([]*sql.DB, 0),
		logger: log.CreateLogger("ConnectionPool"),
		mux:    sync.Mutex{},
		num:    uint(i),
	}

	isActiveDB := false
	for range numOfConnections {
		connection, err := connect()
		if err == nil {
			connectionPool.c = append(connectionPool.c, connection)
			connectionPool.currNum++
			isActiveDB = true
		}
	}
	if !isActiveDB {
		connectionPool.logger.WriteError("CreateConnectionPool(): unable to connect to db")
		return nil
	}

	return connectionPool
}

type ConnectionPool struct {
	c            []*sql.DB
	logger       log.ILogger
	num, currNum uint
	mux          sync.Mutex
}

func (conn *ConnectionPool) GetConnection() *sql.DB {
	conn.mux.Lock()
	defer conn.mux.Unlock()
	n := len(conn.c)

	if conn.c == nil || n == 0 {
		conn.logger.WriteError("GetConnection(): no connections left")
		return nil
	}

	for range n {
		connection := conn.c[0]
		conn.c = conn.c[1:]
		err := connection.Ping()
		if err == nil {
			if conn.num != conn.currNum {
				conn.logger.WriteWarning(fmt.Sprintf("GetConnection(): %d/%d connections with db active", conn.currNum, conn.num))
				go conn.tryReconnect()
			}
			return connection
		} else {
			connection.Close()
			conn.currNum--
		}
	}

	conn.logger.WriteError("GetConnection(): no connections left")
	go conn.tryReconnect()

	return nil
}

func (conn *ConnectionPool) tryReconnect() {
	conn.mux.Lock()
	defer conn.mux.Unlock()

	delta := conn.num - conn.currNum
	for range delta {
		connection, err := connect()
		if err != nil {
			conn.c = append(conn.c, connection)
			conn.currNum++
		} else {
			conn.logger.WriteError(fmt.Sprintf("tryReconnect(): %v", err))
		}
	}
}

func (conn *ConnectionPool) Close() {
	conn.mux.Lock()
	defer conn.mux.Unlock()
	for _, connection := range conn.c {
		connection.Close()
	}
	conn.num = 0
	conn.currNum = 0
}
