package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type TracerouteCLIOptions struct {
	dispHopNumeric bool `default:"false"`
	nQueries       int  `default:"1"`
	dispProbeStat  bool `default:"false"`
}

type ProbeResult struct {
	isSuccess   bool `binding:"required"`
	ip          net.Addr
	timeElapsed time.Duration
}

type AddrResults map[string][]*ProbeResult

type contextKey int

const (
	socketKey contextKey = iota
	dstIpKey
	mutexKey
)

const (
	packetSize = 64
	numMaxHops = 30
)

func NewTracerouteCLIOptions(dispHopsNumeric bool, nqueries int, dispProbeStat bool) *TracerouteCLIOptions {
	return &TracerouteCLIOptions{
		dispHopNumeric: dispHopsNumeric,
		nQueries:       nqueries,
		dispProbeStat:  dispProbeStat,
	}
}

func displayHopStats(success int, total int) {
	rate := float32((total - success)) / float32(total) * 100.0
	fmt.Printf(" -- Successes %d Losses %d Loss Rate %.2f%%", success, total-success, rate)
}

func dispRepeatTime(t time.Duration) {
	fmt.Printf(" %8s", t.Truncate(time.Microsecond))
}

func dispSuccessNonSymbolic(startChar string, ip net.Addr, t time.Duration) {
	fmt.Printf("%s%-10s %8s", startChar, ip, t.Truncate(time.Microsecond))
}

func dispSuccess(startChar string, symbol string, ip net.Addr, t time.Duration) {
	fmt.Printf("%s%-10s (%-10s) %8s", startChar, symbol, ip, t.Truncate(time.Microsecond))
}

func dispFail() {
	fmt.Printf(" %-5s", "*")
}

func parseFlags() *TracerouteCLIOptions {
	var nqueries int
	var isHopNumericDisp, isProbeStatDisp bool

	flag.BoolVar(&isHopNumericDisp, "n", false, "Print hop addresses numerically")
	flag.IntVar(&nqueries, "q", 3, "Set the number of probes per `ttl` to nqueries")
	flag.BoolVar(&isProbeStatDisp, "S", false, "Print a summary of how many probes were not answered for each hop")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "dst\n    the destination IP address\n")
	}

	flag.Parse()
	return NewTracerouteCLIOptions(isHopNumericDisp, nqueries, isProbeStatDisp)
}

func ping(ctx context.Context) *ProbeResult {
	socket := ctx.Value(socketKey).(*icmp.PacketConn)
	dst := ctx.Value(dstIpKey).(*net.IPAddr)

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Data: make([]byte, packetSize),
		},
	}

	wb, err := message.Marshal(nil)

	if err != nil {
		fmt.Printf("Error marshalling message: %s\n", err)
	}

	socket.SetReadDeadline(time.Now().Add(time.Duration(20 * time.Second)))

	var b = make([]byte, 1500)
	isSuccessful := true

	startTime := time.Now()
	if _, err = socket.WriteTo(wb, dst); err != nil {
		isSuccessful = false
		return &ProbeResult{
			isSuccess:   isSuccessful,
			timeElapsed: time.Duration(0),
			ip:          nil,
		}
	}

	_, addr, err := socket.ReadFrom(b)
	elapsedTime := time.Since(startTime)

	if err != nil {
		isSuccessful = false
	}

	return &ProbeResult{
		isSuccess:   isSuccessful,
		timeElapsed: elapsedTime,
		ip:          addr,
	}
}

func traceroute(ctx context.Context, cliOptions *TracerouteCLIOptions) {
	destinationFound := false
	ttl := 1
	socket := ctx.Value(socketKey).(*icmp.PacketConn)
	dstIp := ctx.Value(dstIpKey).(*net.IPAddr)

	for ttl < numMaxHops && !destinationFound {
		var successCnt int = 0
		var probeCount int = cliOptions.nQueries
		socket.IPv4PacketConn().SetTTL(ttl)

		fmt.Printf("%-5d", ttl)
		var prevIp net.Addr = nil

		for i := 0; i < probeCount; i++ {
			result := ping(ctx)

			if result.isSuccess {
				successCnt++
				if prevIp != nil && prevIp.String() == result.ip.String() {
					dispRepeatTime(result.timeElapsed)
				} else {
					var startChar string

					if prevIp == nil {
						startChar = ""
					} else {
						startChar = "\n"
					}

					if !cliOptions.dispHopNumeric {
						domain, err := net.LookupAddr(result.ip.String())

						if err != nil {
							dispSuccess(startChar, result.ip.String(), result.ip, result.timeElapsed)
						} else {
							dispSuccess(startChar, string(domain[0]), result.ip, result.timeElapsed)
						}
					} else {
						dispSuccessNonSymbolic(startChar, result.ip, result.timeElapsed)
					}
				}

				if result.ip.String() == dstIp.String() {
					destinationFound = true
				}
			} else {
				dispFail()
			}
			prevIp = result.ip
		}

		if cliOptions.dispProbeStat {
			displayHopStats(successCnt, probeCount)
		}
		fmt.Println()
		ttl += 1
	}
}

func main() {
	var dst string

	cliOptions := parseFlags()
	ctx := context.Background()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	dst = flag.Arg(0)
	dstIp, err := net.ResolveIPAddr("ip4", dst)

	if err != nil {
		fmt.Println("Invalid IP address format.")
		os.Exit(1)
	}
	ctx = context.WithValue(ctx, dstIpKey, dstIp)

	socket, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		fmt.Print(err)
	}
	defer socket.Close()

	ctx = context.WithValue(ctx, socketKey, socket)
	ctx = context.WithValue(ctx, mutexKey, &sync.Mutex{})

	traceroute(ctx, cliOptions)
}
