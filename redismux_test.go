package redismux

import (
	"testing"
)

func TestMuxer(t *testing.T) {
	mux, err := NewRedisMux("127.0.0.1:6379")
	if err != nil {
		t.Fatalf("Couldn't create redis muxer: %s", err)
	}

	res1, err := mux.Do("set", "foo", 42)
	if err != nil {
		t.Fatalf("Got an error from redis: %s", err)
	}
	t.Log("Got reply %s", res1)

	res2, err := mux.Do("get", "foo")
	if err != nil {
		t.Fatalf("Got an error from redis: %s", err)
	}
	t.Log("Got reply %s", res2)

	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func() {
			for n := 0; n < 50000; n++ {
				_, err := mux.Do("INCRBY", "FOO", 1)
				if err != nil {
					t.Fatalf("Got an error from redis: %s", err)
				}
			}
			ch <- 1
		}()
	}
	for i := 0; i < 5; i++ {
		<-ch
	}
}
