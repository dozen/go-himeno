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

type Manager struct {
	Size      string
	Nodes     []*pb.Node
	StartLock sync.RWMutex
	KickLock  sync.RWMutex
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
	return &res, nil
}

func (mc *Manager) Join(ctx context.Context, in *pb.JoinReq) (*pb.JoinRes, error) {
	//for Worker
	//score と linkspeed を申告して参加
	res := pb.JoinRes{}
	mc.Nodes = append(mc.Nodes, &pb.Node{
		Status:    "ok",
		Address:   in.Addr,
		Score:     in.Score,
		LinkSpeed: in.LinkSpeed,
	})
	return &res, nil
}

func (mc *Manager) Job(ctx context.Context, in *pb.JobReq) (*pb.JobRes, error) {
	//Jobの割り当てが終わるのを待って各ノードにJobを送信
	mc.StartLock.RLock()
	res := pb.JobRes{}
	return &res, nil
}

func (mc *Manager) Kick(ctx context.Context, in *pb.KickReq) (*pb.KickRes, error) {
	//各ノードは他のノードと接続ができ次第 KickReq を送る。
	//全部の KickReq が送られてきたのを確認して KickRes を一斉に返す
	res := pb.KickRes{}

	mc.KickLock.RLock()
	return &res, nil
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
	mng.KickLock.Lock()  //Job生成が終わったら解除

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
