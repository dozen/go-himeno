package main

import (
	"context"
	"fmt"
	pb "github.com/dozen/go-himeno/manager/proto"
	"gopkg.in/alecthomas/kingpin.v2"
	"net"
	"sync"
)

const (
	Protocol = "tcp"
)

var (
	mngAddr = kingpin.Flag("manager", "Manager Host:Port").Short('m').Required().String()
	addr    = kingpin.Flag("listen", "Listen Host:Port").Short('l').Default(":22123").String()
	score   = kingpin.Flag("score", "Score MFLOPS").Short('s').Default("3000").Float64()

	job pb.JobRes
)

func init() {
	kingpin.Parse()

	go NeighborServer()
}

func join(ctx context.Context, c pb.ManagerClient) {
	_, err := c.Join(ctx, &pb.JoinReq{
		Addr:      *addr,
		Score:     *score,
		LinkSpeed: 1000,
	})
	if err != nil {
		panic(err)
	}
}

func getJob(ctx context.Context, c pb.ManagerClient) {
	r, err := c.Job(ctx, &pb.JobReq{
		Addr: *addr,
	})
	if err != nil {
		panic(err)
	}
	job = *r

	startCommunication()
}

func waitKick(ctx context.Context, c pb.ManagerClient) {
	_, err := c.Kick(ctx, &pb.KickReq{
		Addr: *addr,
	})
	if err != nil {
		panic(err)
	}
}

func reportTimes(ctx context.Context, c pb.ManagerClient, times int) int {
	fmt.Println("Suggest Times:", times, "sending...")
	r, err := c.ReportTimes(ctx, &pb.ReportTimesReq{
		Addr:  *addr,
		Times: int64(times),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Send.")
	return int(r.Times)
}

func result(ctx context.Context, c pb.ManagerClient, gosa float32, cpu float64) {
	_, err := c.Result(ctx, &pb.ResultReq{
		Addr: *addr,
		Gosa: gosa,
		Cpu:  cpu,
	})
	if err != nil {
		panic(err)
	}
}

func startCommunication() {
	//ここからクライアント
	clientsWait := sync.WaitGroup{}
	if job.LeftNeighbor != "" {
		clientsWait.Add(1)
		go func() {
			NeighborClient(job.LeftNeighbor, int(job.Left), "left")
			clientsWait.Done()
		}()
	}
	if job.RightNeighbor != "" {
		clientsWait.Add(1)
		go func() {
			NeighborClient(job.RightNeighbor, int(job.Right), "right")
			clientsWait.Done()
		}()
	}

	clientsWait.Wait()
}

func NeighborClient(addr string, index int, dist string) {
	remoteIndex := 0 //お隣さんの [i]
	if dist == "left" {
		remoteIndex = index - 1
	} else {
		remoteIndex = index + 1
	}

	host, err := net.ResolveTCPAddr(Protocol, addr)
	if err != nil {
		fmt.Print(addr, ": ")
		panic(err)
	}

	conn, err := net.DialTCP(Protocol, nil, host)
	if err != nil {
		fmt.Print(addr, ": ")
		panic(err)
	}

	// TODO: ClientHandlerでハンドシェイク的なの済ませてから抜けた方がいい気がしてきた
	go ClientHandler(conn, index, remoteIndex)
}

func ClientHandler(conn *net.TCPConn, local, remote int) {
	//通信で Neighbor に 隣の (remoteの) 配列をもらう処理を書く
}

func NeighborServer() {
	//Serve Neighbor Communication
	src, err := net.ResolveTCPAddr(Protocol, *addr)
	if err != nil {
		fmt.Print(*addr, ": ")
		fmt.Println("ServeNC resovle TCP addr error.")
		panic(err)
	}

	lis, err := net.ListenTCP(Protocol, src)
	if err != nil {
		fmt.Print(*addr, ": ")
		fmt.Println("ServeNC listen TCP error.")
		panic(err)
	}
	defer func() {
		lis.Close()
		fmt.Println(*addr, ": ServerNC Close.")
	}()
	fmt.Print(*addr, ": ")
	fmt.Println("Listen Start Neighbor Server")

	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			fmt.Print(*addr, ": ")
			fmt.Println("ServeNC accept error.")
			fmt.Println(err)
			if conn != nil {
				fmt.Println(conn.RemoteAddr().String())
			}
			continue
		}
		go NCHandler(conn)
	}
}

func NCHandler(conn *net.TCPConn) {
	fmt.Println(conn.LocalAddr().String(), ": connected by", conn.RemoteAddr().String())

	reqN := make([]byte, 4)
	_, err := conn.Read(reqN)
	if err != nil {
		panic(err)
	}

}
