package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
	// Import the generated package
	pb "scistream-project/scistream"
)

func main() {
        // Hardcoded values
        uid := "745cf1d0a3aa11efbec80242ac110003"
        prodLstn := "192.168.10.10,192.168.10.10:5075,192.168.10.10:5076"
        s2cs := "192.168.10.11:5000"

	// Connect to server
	ctx := context.Background()
	conn, err := grpc.Dial(s2cs, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	defer conn.Close()
	
	// Create client
	client := pb.NewControlClient(conn)
	fmt.Printf("Session: %s %s CONS\n", uid, s2cs)
	
	// Use WaitGroup to manage goroutines
	var wg sync.WaitGroup
	
	// Send client request in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		fmt.Println("Sending client request")
		clientResp, err := client.Req(ctx, &pb.Request{
			Uid: uid, 
			Role: "CONS", 
			NumConn: 5, 
			Rate: 10000,
		})
		
		if err != nil {
			fmt.Printf("Client request failed: %v\n", err)
		} else {
			fmt.Printf("Client response: %v\n", clientResp)
		}
	}()
	
	// Short pause before sending Hello
	time.Sleep(500 * time.Millisecond)
	
	// Send hello request with prod listeners
	fmt.Println("Sending hello message")
	prodListeners := strings.Split(prodLstn, ",")
	
	helloResp, err := client.Hello(ctx, &pb.Hello{
		Uid: uid, 
		Role: "PROD", 
		ProdListeners: prodListeners,
	})
	if err != nil {
		fmt.Printf("Hello request failed: %v\n", err)
		return
	}
	
	fmt.Printf("Hello response: %v\n", helloResp)
	
	// Send update request with same listeners
	updateResp, err := client.Update(ctx, &pb.UpdateTargets{
		Uid: uid, 
		RemoteListeners: prodListeners, 
		Role: "CONS",
	})
	
	if err != nil {
		fmt.Printf("Update failed: %v\n", err)
	} else {
		fmt.Printf("Update successful: %v\n", updateResp)
	}
	
	// Wait for client request goroutine to finish (optional)
	// Comment this out if you want the program to exit without waiting
	// for the client request to complete
	wg.Wait()
}
