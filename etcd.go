package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

func main() {
	cli1, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli1.Close()

	cli2, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli2.Close()

	// create two separate sessions for lock competition
	s1, err := concurrency.NewSession(cli1)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()

	m1 := concurrency.NewMutex(s1, "/my-lock/")

	s2, err := concurrency.NewSession(cli2)
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()
	m2 := concurrency.NewMutex(s2, "/my-lock/")

	for i := 0; i < 1000; i++ {
		cli1.Put(context.Background(), fmt.Sprintf("/tasks/v1/claim/%04d", i), "")
	}

	for i := 0; i < 1000; i++ {
		if i%7 == 0 {
			continue
		}
		cli1.Put(context.Background(), fmt.Sprintf("/tasks/v1/claim/%04d/owner", i), "hostname", clientv3.WithLease(s1.Lease()))

	}

	fmt.Println("12")
	stuff := ListUnclaimedTasks(cli2, 12)
	if len(stuff) != 12 {
		fmt.Println("not enough stuff 12")
	}
	fmt.Println(stuff)
	fmt.Println("230")
	stuff = ListUnclaimedTasks(cli2, 230)
	if len(stuff) != 230 {
		fmt.Println("not enough stuff 230", len(stuff))
	}
	fmt.Println(stuff)

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(cli1.Delete(context.Background(), "/tasks/v1/claim/", clientv3.WithPrefix()))

	kv1 := clientv3.NewKV(cli1)
	kv1.Put(context.TODO(), "dude", "guy")

	kv1.Put(context.TODO(), "dude", "what", clientv3.WithLease(s1.Lease()))
	r, err := kv1.Get(context.TODO(), "dude")
	fmt.Println("get dude", r, err)

	kv2 := clientv3.NewKV(cli2)
	kv2.Put(context.TODO(), "man", "hi", clientv3.WithLease(s2.Lease()))
	r, err = kv2.Get(context.TODO(), "man")
	fmt.Println("get man", r, err)

	// acquire lock for s1
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// wait until s1 is locks /my-lock/
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("acquired lock for s2")
		if err := m2.Unlock(context.TODO()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("released lock for s2")
	}()

	s1.Close()
	time.Sleep(10 * time.Second)

	r, err = kv1.Get(context.TODO(), "dude")
	fmt.Println("get dude", r, err)

	r, err = kv2.Get(context.TODO(), "man")
	fmt.Println("get man", r, err)

	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")

	<-m2Locked
	fmt.Println("done")

	r, err = kv1.Get(context.TODO(), "dude")
	fmt.Println("get dude", r, err)

	r, err = kv2.Get(context.TODO(), "man")
	fmt.Println("get man", r, err)

}

func ListUnclaimedTasks(cli *clientv3.Client, count int) []string {

	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithFromKey(),
		clientv3.WithLimit(int64(count)),
	}

	r, err := cli.Get(context.Background(), "/tasks/v1/claim/", opts...)
	if err != nil {
		log.Fatal(err)
	}

	// a simple way of guaranteeing uniqueness
	availableLeases := map[string]struct{}{}

	num := 0
	for err == nil {
		num++

		for _, kv := range r.Kvs {
			if bytes.Contains(kv.Key, []byte("owner")) {
				delete(availableLeases, string(bytes.Trim(bytes.Trim(kv.Key, "/tasks/v1/claim/"), "/owner")))
			} else {
				if len(availableLeases) >= count {
					fmt.Println("number reached")
					goto TALLY
				}
				availableLeases[string(bytes.Trim(bytes.Trim(kv.Key, "/tasks/v1/claim/"), "/owner"))] = struct{}{}
			}

		}

		if len(r.Kvs) < count {
			goto TALLY
		}

		lastKey := string(r.Kvs[len(r.Kvs)-1].Key)
		r, err = cli.Get(context.Background(), lastKey, opts...)
		if err != nil {
			log.Fatal(err)
		}

	}

TALLY:

	resp := []string{}
	for key, _ := range availableLeases {
		resp = append(resp, key)
	}

	return resp
}
