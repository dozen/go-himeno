package main

import (
	"context"
	"fmt"
	"github.com/dozen/go-himeno/manager"
	pb "github.com/dozen/go-himeno/manager/proto"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	addr    = kingpin.Flag("addr", "Server Address").Default("127.0.0.1:" + manager.MngPort).String()
	command = kingpin.Arg("command", "Command").Required().String()
)

func main() {
	kingpin.Parse()

	c, closer := manager.ManagerClient(*addr)
	defer closer()

	switch *command {
	case "stats":
		stats(c)
	case "start":
		start(c)
	}
}

func stats(c pb.ManagerClient) {
	ctx := context.Background()
	r, err := c.Stats(ctx, &pb.StatsReq{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", r.NodeList)
}

func start(c pb.ManagerClient) {
	ctx := context.Background()
	_, err := c.Start(ctx, &pb.StartReq{})
	if err != nil {
		panic(err)
	}

	fmt.Println("start.")
}
