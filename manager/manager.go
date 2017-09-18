package manager

import (
	"context"
	"fmt"
	pb "github.com/dozen/go-himeno/manager/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

const (
	Protocol = "tcp"
	MngPort  = "22122"
	MngAddr  = "0.0.0.0:" + MngPort
)

var (
	Size2MIMAX = map[string]int{
		"SSMALL": 33,
		"SMALL":  65,
		"MIDDLE": 129,
		"LARGE":  257,
		"ELARGE": 513,
	}
)

type Manager struct {
	Size      string
	Nodes     []*pb.Node
	NodesLock sync.RWMutex
	StartLock sync.RWMutex
	KickCount sync.WaitGroup
}

func (mc *Manager) Stats(ctx context.Context, in *pb.StatsReq) (*pb.StatsRes, error) {
	//for CLI
	//ノードの一覧などを返す
	res := pb.StatsRes{mc.Nodes}
	return &res, nil
}

func (mc *Manager) Start(ctx context.Context, in *pb.StartReq) (*pb.StartRes, error) {
	//for CLI
	//Joinを締め切ってJobの割り当てを行う
	res := pb.StartRes{}
	//jobの振り分けを実装
	mc.setJob()
	mc.StartLock.Unlock() //Job送ってもいいようにする
	return &res, nil
}

func (mc *Manager) Join(ctx context.Context, in *pb.JoinReq) (*pb.JoinRes, error) {
	//for Worker
	//score と linkspeed を申告して参加
	res := pb.JoinRes{}

	mc.NodesLock.Lock()
	mc.Nodes = append(mc.Nodes, &pb.Node{
		Status:    "ok",
		Address:   in.Addr,
		Score:     in.Score,
		LinkSpeed: in.LinkSpeed,
	})
	mc.NodesLock.Unlock()
	mc.KickCount.Add(1)
	in.Reset()
	return &res, nil
}

func (mc *Manager) Job(ctx context.Context, in *pb.JobReq) (*pb.JobRes, error) {
	//Jobの割り当てが終わるのを待って各ノードにJobを送信
	mc.StartLock.RLock()
	defer mc.StartLock.RUnlock()

	var res *pb.JobRes
	for _, node := range mc.Nodes {
		if node.Address == in.Addr {
			res = node.Job
		}
	}
	return res, nil
}

func (mc *Manager) Kick(ctx context.Context, in *pb.KickReq) (*pb.KickRes, error) {
	//各ノードは他のノードと接続ができ次第 KickReq を送る。
	//全部の KickReq が送られてきたのを確認して KickRes を一斉に返す
	res := pb.KickRes{}

	fmt.Println(in.Addr, "is ready.")
	mc.KickCount.Done()
	mc.KickCount.Wait()
	return &res, nil
}

func (mc *Manager) setJob() {
	totalScore := 0.0
	for _, node := range mc.Nodes {
		totalScore += node.Score
	}

	loads := []int{}
	totalLoad := 0
	for _, node := range mc.Nodes {
		load := int(node.Score/totalScore * float64(Size2MIMAX[mc.Size]))
		loads = append(loads, load)
		totalLoad += load
	}

	//totalLoadがMIMAXより足りない時がある
	//同じ数になるまで揃える処理をする
	for {
		if totalLoad < Size2MIMAX[mc.Size] {
			val := 0
			index := 0
			for i, load := range loads {
				if val < load {
					val = load
					index = i
				}
			}
			loads[index]++
			totalLoad++
		} else {
			break
		}
	}

	//Job作成
	left := 0
	right := 0
	leftNeighbor := ""
	rightNeighbor := ""
	for i, load := range loads {
		left = right
		right = left + load
		leftNeighbor = rightNeighbor
		if i < len(mc.Nodes)-1 {
			rightNeighbor = mc.Nodes[i+1].Address
		} else {
			rightNeighbor = ""
		}
		mc.Nodes[i].Job = &pb.JobRes{
			Size:          mc.Size,
			Left:          int64(left),
			Right:         int64(right),
			LeftNeighbor:  leftNeighbor,
			RightNeighbor: rightNeighbor,
		}
	}
}

func ServeManager(protocol, addr, size string, s *grpc.Server) {
	if protocol == "" {
		protocol = Protocol
	}
	if addr == "" {
		addr = MngAddr
	}
	serverAddr, err := net.ResolveTCPAddr(protocol, addr)
	if err != nil {
		panic(err)
	}
	lis, err := net.ListenTCP(protocol, serverAddr)
	if err != nil {
		panic(err)
	}

	mng := Manager{
		Size: size,
	}
	mng.StartLock.Lock() //Join終わったら解除

	pb.RegisterManagerServer(s, &mng)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		fmt.Errorf("%#v\n", err)
	}
}

func ManagerClient(addr string) (pb.ManagerClient, func() error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := pb.NewManagerClient(conn)

	return c, conn.Close
}
