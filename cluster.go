package clickhouse

import (
	"math/rand"
	"sync"
)

type PingErrorFunc func(*Conn)

type Cluster struct {
	conn   []*Conn
	active []*Conn
	fail   PingErrorFunc
	mx     sync.Mutex
}

func NewCluster(conn ...*Conn) *Cluster {
	return &Cluster{
		conn: conn,
	}
}

func (c *Cluster) OnPingError(f PingErrorFunc) {
	c.fail = f
}

func (c *Cluster) ActiveConn() *Conn {
	c.mx.Lock()
	defer c.mx.Unlock()
	l := len(c.active)
	if l < 1 {
		return nil
	}
	return c.active[rand.Intn(l)]
}

func (c *Cluster) Ping() {
	var (
		err error
		res []*Conn
	)

	for _, conn := range c.conn {
		err = conn.Ping()
		if err == nil {
			res = append(res, conn)
		} else {
			c.fail(conn)
		}
	}

	c.mx.Lock()
	c.active = res
	c.mx.Unlock()
}