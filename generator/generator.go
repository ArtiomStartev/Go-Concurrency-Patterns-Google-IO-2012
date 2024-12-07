package generator

import (
	"fmt"
	"math/rand"
	"time"
)

// Generator: function that returns a channel
// Channels are first-class values, just like strings or integers

func Init() {
	// Our boring function returns a channel that let us communicate with the boring service it provides
	сh := boring("boring!") // Function returning a channel
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-сh) // Read from the channel
	}

	// We can have more instances of the service
	joe := boring("Joe")
	ann := boring("Ann")
	for i := 0; i < 5; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}

	// We can also have a fanIn function that multiplexes the channels
	// joe channel does not block ann channel and vice versa anymore
	c := fanIn(joe, ann)
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

	fmt.Println("You're boring; I'm leaving.")
}

func boring(msg string) <-chan string { // Returns receive-only channel of strings
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller
}
