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

	src = kingpin.Flag("src", "Client Address").Default("0.0.0.0").String()
	score = kingpin.Flag("score", "Client Score").Default("3000").Float64()
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
	case "join":
		join(c)
	}
}

func stats(c pb.ManagerClient) {
	ctx := context.Background()
	r, err := c.Stats(ctx, &pb.StatsReq{})
	if err != nil {
		panic(err)
	}

	for _, node := range r.NodeList {
		fmt.Printf("%v { score:%v } : %+v\n", node.Address, node.Score, *(node.Job))
	}
}

func start(c pb.ManagerClient) {
	ctx := context.Background()
	_, err := c.Start(ctx, &pb.StartReq{})
	if err != nil {
		panic(err)
	}

	fmt.Println("start.")
}

func join(c pb.ManagerClient) {
	ctx := context.Background()
	_, err := c.Join(ctx, &pb.JoinReq{
		Addr: *src,
		Score: *score,
		LinkSpeed: 1000,
	})
	if err != nil {
		panic(err)
	}
}