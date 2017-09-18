package main

import (
	"github.com/dozen/go-himeno/manager"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	size = kingpin.Flag("size", "himeno size").Default("LARGE").String()
	addr = kingpin.Flag("addr", "server addr").Default(manager.MngAddr).String()
)

func main() {
	kingpin.Parse()

	s := grpc.NewServer()
	manager.ServeManager("", *addr, *size, s)
}
