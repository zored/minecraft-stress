package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/ztrue/shutdown"
	"log"
	"sync"
	"time"
)

type (
	pingOutput struct {
		Duration time.Duration
		List     *json.RawMessage
	}
)

var (
	host       = flag.String("host", "localhost", "server host")
	port       = flag.Int("port", 25565, "server port")
	bots       = flag.Int("bots", 20, "number of bots")
	ping       = flag.Bool("ping", false, "just ping server")
	nickPrefix = flag.String("nick", "man", "bot nick prefix")
)

func init() {
	flag.Parse()
}

func main() {
	if *ping {
		doPing()
	} else {
		doStressTest()
	}
}

func doStressTest() {
	clients := make([]*bot.Client, 0)
	clientMutex := &sync.Mutex{}

	shutdown.Add(func() {
		clientMutex.Lock()
		for _, client := range clients {
			err := client.Disconnect()
			if err != nil {
				log.Fatal(err)
			}
		}
		clientMutex.Unlock()
	})

	wg := sync.WaitGroup{}
	for number := 1; number <= *bots; number++ {
		time.Sleep(time.Second)
		wg.Add(1)
		go func(number int) {
			client := bot.NewClient()
			clientMutex.Lock()
			clients = append(clients, client)
			clientMutex.Unlock()

			nick := fmt.Sprintf("%s_%d", *nickPrefix, number)
			client.Name = nick
			h := *host
			p := *port
			err := client.JoinServer(h, p)
			panicIfError(err)
			fmt.Printf("'%s' joined \n", nick)

			client.Events.GameStart = func() error {
				return client.Dig(0, 0, 0, 0, 0)
			}

			err = client.HandleGame()
			panicIfError(err)

			wg.Done()
		}(number)
	}

	wg.Wait()
}

func doPing() {
	list, duration, err := bot.PingAndList(*host, *port)
	panicIfError(err)
	rawList := json.RawMessage(list)

	output, err := json.Marshal(&pingOutput{
		Duration: duration,
		List:     &rawList,
	})
	panicIfError(err)
	fmt.Printf("%s", output)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
