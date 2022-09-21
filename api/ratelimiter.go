package api

import (
	"golang.org/x/time/rate"
	"net"
	"sync"
	"time"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type visitors struct {
	visitors map[string]*visitor
	lock     sync.Mutex
}

func (v *visitors) addVisitor(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(5, 500)
	v.lock.Lock()
	v.visitors[ip] = &visitor{limiter, time.Now()}
	v.lock.Unlock()
	return limiter
}

func (v *visitors) getVisitor(ip string) *rate.Limiter {
	v.lock.Lock()
	visitor, exists := v.visitors[ip]
	if !exists {
		v.lock.Unlock()
		return v.addVisitor(ip)
	}

	visitor.lastSeen = time.Now()
	v.lock.Unlock()
	return visitor.limiter
}

func (v *visitors) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		v.lock.Lock()
		for ip, visitor := range v.visitors {
			if time.Now().Sub(visitor.lastSeen) > 3*time.Minute {
				delete(v.visitors, ip)
			}
		}
		v.lock.Unlock()
	}
}

func (v *visitors) isAllowed(hostPort string) bool {
	ip, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		return false
	}
	limiter := v.getVisitor(ip)
	return limiter.Allow()
}

func newVisitors() *visitors {
	visitors := &visitors{
		visitors: make(map[string]*visitor),
	}
	return visitors
}
