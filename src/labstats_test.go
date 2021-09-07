package main

import (
	"testing"
)

func TestMachineName(t *testing.T) {
  c := NewStatusClient(&StatusOptions{})
  if c.options.Name == "TestMachine" {
    t.Logf("Get Machine Name success!")
  } else {
    t.Errorf("expected Name = TestMachine. Found %v", c.options.Name)
  }
}


// gRPC information
const (
  address = "localhost:50051"
)

func TestCallStats(t *testing.T) {
  s := NewStatusServer(&StatusOptions{})
  c := NewStatusClient(&StatusOptions{
    ServerAddr: "localhost:50051",
  })
  go s.Start()
  r_val := c.CallStats()
  s.grpcServer.Stop()
  if r_val != nil {
    t.Errorf("Unexpected Error Found %v", r_val)
  } else {
    t.Logf("gRPC Stats success!")
  }
}

