package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type PingCLIOptions struct {
	Count      int `default:"-1"`
	Wait       int `default:"1"`
	PacketSize int `default:"0"`
	Timeout    int `default:"3"`
}

type PingStats struct {
	Successes  int
	Failures   int
	PacketSize int
	TTL        int
	Seq        int
}

type contextKey int

const (
	socketKey contextKey = iota
	dstKey
	timeoutDurationKey
	ttlKey
	waitDurationKey
	packetSizeKey
)

func NewPingCLIOptions(count int, wait int, ps int, timeout int) *PingCLIOptions {
	return &PingCLIOptions{
		Count:      count,
		Wait:       wait,
		PacketSize: ps,
		Timeout:    timeout,
	}
}

func parseFlags() *PingCLIOptions {
	var count, wait, packetSize, timeout int
	flag.IntVar(&count, "c", -1, "Number of pings to send")
	flag.IntVar(&wait, "i", 1, "Time to wait between pings in seconds")
	flag.IntVar(&packetSize, "s", 8, "Packet size to deliver")
	flag.IntVar(&timeout, "t", 3, "Timeout before ping exits in seconds")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "dst\n    the destination IP address\n")
	}

	flag.Parse()
	return NewPingCLIOptions(count, wait, packetSize, timeout)
}

func displayStats(stats *PingStats, totalTime time.Duration) {
	aveTime := totalTime / time.Duration(stats.Seq)
	// Add one to prevent from divide by zero
	packetLossRate := (1 - (float32(stats.Successes) / float32(stats.Seq))) * 100.0

	fmt.Printf("\n---------- Ping Statistics ----------\n")
	fmt.Printf("Average %s, Success: %d, Failure: %d, Paket Loss %.2f %%\n",
		aveTime.String(), stats.Successes, stats.Failures, packetLossRate)
}

func Ping(ctx context.Context, stats *PingStats) time.Duration {
	icmpPacket := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  stats.Seq,
			Data: make([]byte, stats.PacketSize-8),
		},
	}
	wb, err := icmpPacket.Marshal(nil)

	if err != nil {
		fmt.Print(err)
	}
	socket := ctx.Value(socketKey).(*icmp.PacketConn)
	timeoutDuration := ctx.Value(timeoutDurationKey).(time.Duration)
	dst := ctx.Value(dstKey).(net.Addr)

	var msg string
	startTime := time.Now()
	// Set deadline to read the incoming message
	socket.SetReadDeadline(time.Now().Add(time.Duration(timeoutDuration)))
	// Send ping to client
	if _, err = socket.WriteTo(wb, dst); err != nil {
		msg = fmt.Sprintf("%s", err)
	}

	var b []byte
	_, _, err = socket.ReadFrom(b)
	elapsedTime := time.Since(startTime)

	if err != nil {
		msg = fmt.Sprintf("%s", err)
	}
	if msg != "" {
		stats.Failures++
		fmt.Printf("%s ttl=%d icmp_seq=%d: %s\n",
			ctx.Value(dstKey), stats.TTL, stats.Seq, msg)
	} else {
		stats.Successes++
		fmt.Printf("%d bytes transferred from %s ttl=%d icmp_seq=%d roundtrip time %s\n",
			stats.PacketSize, ctx.Value(dstKey), stats.TTL, stats.Seq, elapsedTime.Truncate(time.Microsecond))
	}

	return elapsedTime
}

func main() {
	var dst *net.IPAddr
	timeTotal := time.Duration(0) // used to calculate statistics
	cliOptions := parseFlags()
	ctx := context.Background()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	dst, err := net.ResolveIPAddr("ip4", os.Args[flag.NFlag()+2])

	if err != nil {
		fmt.Println("Invalid IP address format.")
		os.Exit(1)
	}
	ctx = context.WithValue(ctx, dstKey, dst)

	socket, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		fmt.Print(err)
	}
	ctx = context.WithValue(ctx, socketKey, socket)
	defer socket.Close()

	ttl, err := socket.IPv4PacketConn().TTL()

	if err != nil {
		fmt.Print(err)
	}
	stats := &PingStats{
		PacketSize: cliOptions.PacketSize,
		TTL:        ttl,
		Seq:        0,
		Successes:  0,
		Failures:   0,
	}
	timeoutDuration := time.Duration(cliOptions.Timeout) * time.Second
	waitDuration := time.Duration(cliOptions.Wait) * time.Second

	ctx = context.WithValue(ctx, timeoutDurationKey, timeoutDuration)
	ctx = context.WithValue(ctx, waitDurationKey, waitDuration)

	if cliOptions.Count == -1 {
		for {
			timeTotal += Ping(ctx, stats)
			// Sleep until next ping
			time.Sleep(waitDuration)
			stats.Seq++
		}
	} else {
		count := 0
		for count < cliOptions.Count {
			timeTotal += Ping(ctx, stats)
			// Sleep until next ping
			time.Sleep(waitDuration)
			count++
			stats.Seq++
		}
	}

	displayStats(stats, timeTotal)
}
