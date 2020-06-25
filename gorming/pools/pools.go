package main

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/echo-health/training-materials/gorming/connection"
)

func main() {
	rand.Seed(time.Now().Unix())

	// if you check connection.GetConnection, you'll see that there's a maximum number of open connections we can have. Let's see
	// what happens if you try to get more than that.
	i := 0
	for i < 11 {
		conn, err := connection.GetConnection()
		if err != nil {
			panic(err)
		}
		open := make(chan bool)
		go func() {
			// Every time you do this, you're opening a connection, 'taking it off the pool'. When the pool is exhausted, if we
			// try to open a new connection in this line, it'll wait until there's a new connection in the pool
			fmt.Println("--------------------")
			fmt.Printf("Waiting to acquire connection #%d\n", i)
			tx := conn.Begin()
			fmt.Printf("Got connection #%d\n", i)
			open <- true
			time.Sleep(5 * time.Second)
			tx.Rollback()
		}()
		<-open
		i++
	}
}
