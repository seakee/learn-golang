package main

import (
	"fmt"
	"net"
	"sync"
	"time"
	"unsafe"
)

func main() {
	tcpScanByGoroutineWithChannel("127.0.0.1", 1, 65535)
}

func handleWorker(ip string, ports chan int, wg *sync.WaitGroup) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", ip, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// fmt.Printf("[info] %s Close \n", address)
			wg.Done()
			continue
		}
		conn.Close()
		fmt.Printf("[info] %s Open \n", address)
		wg.Done()
	}
}

func tcpScanByGoroutineWithChannel(ip string, portStart int, portEnd int) {
	start := time.Now()

	// 参数校验
	isok := verifyParam(ip, portStart, portEnd)
	if isok == false {
		fmt.Printf("[Exit]\n")
	}

	ports := make(chan int, 100)
	var wg sync.WaitGroup
	for i := 0; i < cap(ports); i++ {
		go handleWorker(ip, ports, &wg)
	}

	for i := portStart; i <= portEnd; i++ {
		wg.Add(1)
		ports <- i
	}

	wg.Wait()
	close(ports)

	cost := time.Since(start)
	fmt.Printf("[tcpScanByGoroutineWithChannel] cost %s second \n", cost)
}

func verifyParam(ip string, portStart int, portEnd int) bool {
	netip := net.ParseIP(ip)
	if netip == nil {
		fmt.Println("[Error] ip type is must net.ip")
		return false
	}
	fmt.Printf("[Info] ip=%s | ip type is: %T | ip size is: %d \n", netip, netip, unsafe.Sizeof(netip))

	if portStart < 1 || portEnd > 65535 {
		fmt.Println("[Error] port is must in the range of 1~65535")
		return false
	}
	fmt.Printf("[Info] port start:%d end:%d \n", portStart, portEnd)

	return true
}
