package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string
type Search func(query string) Result

var (
	Web1   = fakeSearch("web1-replica")
	Web2   = fakeSearch("web2-replica")
	Image1 = fakeSearch("image1-replica")
	Image2 = fakeSearch("image2-replica")
	Video1 = fakeSearch("video1-replica")
	Video2 = fakeSearch("video2-replica")
)

func main() {
	start := time.Now()
	results := googleSearch3("golang")
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func googleSearch1(query string) (results []Result) {
	c := make(chan Result)

	// Run the Web, Image and Video searches concurrently, and wait for all results
	go func() { c <- Web1(query) }()
	go func() { c <- Image1(query) }()
	go func() { c <- Video1(query) }()

	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}

	return
}

func googleSearch2(query string) (results []Result) {
	c := make(chan Result)

	// Run the Web, Image and Video searches concurrently
	go func() { c <- Web1(query) }()
	go func() { c <- Image1(query) }()
	go func() { c <- Video1(query) }()

	// Don't wait for slow servers
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}

	return
}

// Replicate the servers. Send requests to multiple replicas and use the first response
// This is a form of replication. It's a common pattern for distributed systems
// This way we get the fastest response

func First(query string, replicas ...Search) Result {
	c := make(chan Result)

	for i := range replicas {
		go func(j int) {
			c <- replicas[j](query)
		}(i)
	}

	return <-c
}

func googleSearch3(query string) (results []Result) {
	c := make(chan Result)

	// Run the Web, Image and Video searches concurrently
	go func() { c <- First(query, Web1, Web2) }()
	go func() { c <- First(query, Image1, Image2) }()
	go func() { c <- First(query, Video1, Video2) }()

	// Don't wait for slow servers.
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}

	return
}

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}
