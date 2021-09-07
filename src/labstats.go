package main

import (
    "context"
    // "encoding/json"
    "log"
    "net"
    "os"
    "time"

    "github.com/go-redis/redis"
    "google.golang.org/grpc"
    pb "protos"
)

type StatusOptions struct {
  // port to bind server to
  BindPort string
  // address to connect to grpc server
  ServerAddr string
  // address to connect to redis
  RedisAddr string
  // Machine Name
  Name string
}

type StatusServer struct {
  grpcServer *grpc.Server
  redisClient *redis.Client
  options *StatusOptions
}

type StatusClient struct {
  options *StatusOptions
}

func (o *StatusOptions) init() {
  if o.BindPort == "" {
    o.BindPort = ":50051"
  }
  if o.ServerAddr == "" {
    o.ServerAddr = "10.10.199.14:50051"
  }
  if o.RedisAddr == "" {
    o.RedisAddr = "10.10.199.14:6379"
  }
  if o.Name == "" {
    o.Name = "TestMachine"
  }
}

func NewStatusServer(opt *StatusOptions) *StatusServer {
  opt.init()
  s := StatusServer{}
  s.options = opt
  return &s
}

func NewStatusClient(opt *StatusOptions) *StatusClient {
  opt.init()
  c := StatusClient{}
  c.options = opt
  return &c
}

func (s *StatusServer) ConnectRedis() {
  s.redisClient = redis.NewClient(&redis.Options{
    Addr: s.options.RedisAddr,
    Password: "",
    DB: 0,
  })
}

type server struct {
  pb.UnimplementedLabStatsServer
}

func (s *StatusServer) Start() {
  lis, err := net.Listen("tcp", s.options.BindPort)
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }
  s.grpcServer = grpc.NewServer()
  pb.RegisterLabStatsServer(s.grpcServer, &server{})
  s.ConnectRedis()
  pong, err := s.redisClient.Ping(s.redisClient.Context()).Result()
  log.Printf("%v, %v",pong, err)
  if err := s.grpcServer.Serve(lis); err != nil {
    log.Fatalf("Failed to start grpc server: %v", err)
  }
}

func (s *server) Stats (ctx context.Context, in *pb.MachineState) (*pb.Empty, error) {
  log.Printf("Stats called by %v", in.GetName())
  return &pb.Empty{}, nil
}

func (c *StatusClient) CallStats() (error){
  // Set up a connection to the server.
  conn, err := grpc.Dial(c.options.ServerAddr, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  grpcClient := pb.NewLabStatsClient(conn)

  // Contact the server and print out its response.
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  _, err = grpcClient.Stats(ctx, &pb.MachineState{Name: c.options.Name})
  if err != nil {
    log.Fatalf("could not call Stats: %v", err)
  }
  return err
}

func main(){
  flag := os.Args[1]
  if flag == "s" {
    s := NewStatusServer(&StatusOptions{})
    s.Start()
  } else if (flag == "c") {
    c := NewStatusClient(&StatusOptions{})
    c.CallStats()
  } else {
    log.Fatalf("Failed to recognize command line option")
  }
}
