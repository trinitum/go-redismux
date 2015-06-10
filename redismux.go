package redismux

import (
	"github.com/garyburd/redigo/redis"
)

type response struct {
	resp interface{}
	err  error
}

type rchan chan response

type request struct {
	command string
	args    []interface{}
	rc      rchan
}

type RedisMux struct {
	c chan request
}

func NewRedisMux(address string) (*RedisMux, error) {
	c := make(chan request)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	go func() {
		rc := make(chan rchan)
		go func() {
			for ch := range rc {
				resp, rerr := conn.Receive()
				ch <- response{resp, rerr}
			}
		}()
		for req := range c {
			serr := conn.Send(req.command, req.args...)
			conn.Flush()
			if serr != nil {
				close(c)
				req.rc <- response{nil, err}
				close(rc)
				return
			}
			rc <- req.rc
		}
	}()
	return &RedisMux{c}, nil
}

func (mux *RedisMux) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := make(chan response)
	mux.c <- request{commandName, args, c}
	resp := <-c
	return resp.resp, resp.err
}
