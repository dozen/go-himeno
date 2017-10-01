package main

import (
	"context"
	"fmt"
	pb "github.com/dozen/go-himeno/manager/proto"
	"gopkg.in/alecthomas/kingpin.v2"
	"net"
	"sync"
	"unsafe"
)

const (
	Protocol = "tcp"

	MsgEnd = 9
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
	go ClientHandler(conn, index, remoteIndex, dist)
}

func ClientHandler(conn *net.TCPConn, local, remote int, dist string) {
	//通信で Neighbor に 隣の (remoteの) 配列をもらう処理を書く
	sig := make(chan byte)
	endSig := make(chan struct{})
	if dist == "left" {
		//左側と通信
		sig = leftChan
		endSig = leftDoneChan
		fmt.Println("左側")
		fmt.Println(job.LeftNeighbor)
		conn.Write([]byte(job.LeftNeighbor))
	} else {
		fmt.Println("右側")
		fmt.Println(job.RightNeighbor)
		conn.Write([]byte(job.RightNeighbor))
		//右側と通信
		sig = rightChan
		endSig = rightDoneChan
	}

	for {
		b := <-sig
		_, err := conn.Write([]byte{b})
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			continue
		}

		total := 0
		payload := make([]byte, payloadSize)
		parted := make([]byte, payloadSize)
		for {
			parted = make([]byte, payloadSize)
			part, err := conn.Read(parted)
			total += part
			payload = append(payload, parted[:part]...)
			if err != nil && err.Error() != "EOF" {
				fmt.Println("Client:", err)
				return
			}
			if err != nil && err.Error() == "EOF" || total >= payloadSize {
				break
			}
		}
		mDeserialize(remote, payload)
		endSig <- struct{}{}
	}
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
	defer conn.Close()
	fmt.Println(conn.LocalAddr().String(), ": connected by", conn.RemoteAddr().String())

	index := 0
	b := make([]byte, 128)
	bLen, err := conn.Read(b)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	if string(b[:bLen]) == job.RightNeighbor {
		fmt.Println("右側のクライアントが接続してきた")
		index = int(job.Right) + 1
	} else {
		fmt.Println("左側のクライアントが接続してきた")
		index = int(job.Left)
	}
	for {
		b = make([]byte, 1)
		if _, err := conn.Read(b); err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			return
		}
		if b[0] == byte(MsgEnd) {
			return
		}

		if n, err := conn.Write(mSerialize(index)); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("Server: ", conn.RemoteAddr().String(), " send", n, "bytes.")
		}
	}
}

func mSerialize(index int) []byte {
	// math             math.Float32bits
	// encoding/binary  binary.PutUint32
	//多分wrk2でいいかな
	bjmax := jmax - 1
	bkmax := kmax - 1
	j := 1
	k := 1

	b := make([]byte, payloadSize)
	shift := 0
	for ; j < bjmax; j++ {
		for ; k < bkmax; k++ {
			v := *(*uint32)(unsafe.Pointer(&wrk2[index][j][k]))
			b[shift+0] = byte(v >> 24)
			b[shift+1] = byte(v >> 16)
			b[shift+2] = byte(v >> 8)
			b[shift+3] = byte(v)
			shift += 4
		}
	}
	return b
}

func mDeserialize(index int, b []byte) {
	// math             math.Float32frombits
	// encoding/binary  binary.Uint32
	//多分wrk2でいいかな
	bjmax := jmax - 1
	bkmax := kmax - 1
	j := 1
	k := 1

	var v uint32
	shift := 0
	for ; j < bjmax; j++ {
		for ; k < bkmax; k++ {
			v = uint32(b[shift+3]) |
				uint32(b[shift+2])<<8 |
				uint32(b[shift+1])<<16 |
				uint32(b[shift+0])<<24
			wrk2[index][j][k] = *(*float32)(unsafe.Pointer(&v))
			shift += 4
		}
	}
}
