package main

import (
	"context"
	"fmt"
	"github.com/dozen/go-himeno/manager"
	pb "github.com/dozen/go-himeno/manager/proto"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	command = kingpin.Arg("command", "Command").Required().String()
	host    = kingpin.Flag("host", "Server Address").Default("127.0.0.1:" + manager.MngPort).String()
	size    = kingpin.Flag("size", "JOB SIZE").Default("LARGE").String()
	src     = kingpin.Flag("src", "Client src Address").Default("0.0.0.0:22123").String()
	score   = kingpin.Flag("score", "Client Score").Default("3000").Float64()

	listen = kingpin.Flag("listen", "Listen Address").Default(manager.MngAddr).String()
)

func main() {
	kingpin.Parse()

	switch *command {
	case "manager":
		serve()
	case "stats":
		c, closer := manager.ManagerClient(*host)
		defer closer()
		stats(c)
	case "start":
		c, closer := manager.ManagerClient(*host)
		defer closer()
		start(c)
	case "join":
		c, closer := manager.ManagerClient(*host)
		defer closer()
		join(c)
	case "kill":
		c, closer := manager.ManagerClient(*host)
		defer closer()
		kill(c)
	default:
		kingpin.Usage()
	}
}

func serve() {
	s := grpc.NewServer()
	manager.ServeManager("", *listen, *size, s)
	defer s.Stop()
}

func stats(c pb.ManagerClient) {
	ctx := context.Background()
	r, err := c.Stats(ctx, &pb.StatsReq{})
	if err != nil {
		panic(err)
	}

	pp.Println(r.NodeList)
	//for _, node := range r.NodeList {
	//	fmt.Printf("%v { score:%v } : %+v\n", node.Address, node.Score, node.Job)
	//}
}

func start(c pb.ManagerClient) {
	c, closer := manager.ManagerClient(*host)
	defer closer()
	ctx := context.Background()
	_, err := c.Start(ctx, &pb.StartReq{})
	if err != nil {
		panic(err)
	}

	fmt.Println("start.")
}

func join(c pb.ManagerClient) {
	c, closer := manager.ManagerClient(*host)
	defer closer()
	ctx := context.Background()
	_, err := c.Join(ctx, &pb.JoinReq{
		Addr:      *src,
		Score:     *score,
		LinkSpeed: 1000,
	})
	if err != nil {
		panic(err)
	}
}

func kill(c pb.ManagerClient) {
	c, closer := manager.ManagerClient(*host)
	defer closer()
	ctx := context.Background()
	_, err := c.Kill(ctx, &pb.KillReq{})
	if err != nil {
		panic(err)
	}
}
