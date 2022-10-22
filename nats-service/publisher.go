package main

import (
	"bufio"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	sp, err := stan.Connect("test-cluster", "publisher", stan.NatsURL("nats://localhost:4222"))

	if err != nil {
		logrus.Fatalln(err)
	}

	file, err := os.OpenFile("materials/model.json", os.O_RDONLY, 0600)

	if err != nil {
		logrus.Fatalln(err)
	}
	defer file.Close()

	var str string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		str += scanner.Text()
	}

	if scanner.Err() != nil {
		logrus.Fatalln(scanner.Err())
	}

	err = sp.Publish("WB", []byte(str))

	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Infoln("Done")
}
