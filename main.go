package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultNTPServer = "time.google.com"
const defaultNTPPort = 123

func main() {
	os.Exit(_main(os.Args))
}

func _main(args []string) int {
	srv, p := defaultNTPServer, defaultNTPPort
	var jsonOut bool
	fs := flag.NewFlagSet("ntp", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `USAGE: %s [NTP Host] [options...]

OPTIONS:
	--json  Enable JSON Output

`, filepath.Base(os.Args[0]))
	}
	fs.BoolVar(&jsonOut, "json", false, "enable json output")
	if len(args) > 1 && strings.IndexByte(args[1], '-') != 0 {
		if strings.IndexByte(args[1], ':') >= 0 {
			var err error
			var p1 string
			srv, p1, err = net.SplitHostPort(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
				return 1
			}
			var p2 uint64
			p2, err = strconv.ParseUint(p1, 10, 16)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
				return 1
			}
			p = int(p2)
		} else {
			srv = args[1]
		}
		// remove the host from the args so fs.Parse can work properly
		args = append([]string{args[0]}, args[2:]...)
	}
	if err := fs.Parse(args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
		return 1
	}
	nt, err := getNetworkTime(srv, p)
	lt := time.Now().UTC()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}
	if jsonOut {
		d, _ := json.Marshal(map[string]interface{}{
			"server":       srv,
			"port":         p,
			"local_time":   lt.Format(time.RFC3339Nano),
			"network_time": nt.Format(time.RFC3339Nano),
		})
		os.Stdout.Write(d)
		return 0
	}
	fmt.Printf("(server: %s, port: %d)\n", srv, p)
	fmt.Printf("local time:   %s\n", lt.Format(time.RFC3339Nano))
	fmt.Printf("network time: %s\n", nt.Format(time.RFC3339Nano))
	return 0
}

// copy-pasted from https://github.com/bt51/ntpclient/blob/3045f71e2530290e28162bbd5bf931ff35f04658/ntpclient.go#L16
func getNetworkTime(srv string, p int) (*time.Time, error) {
	var second, fraction uint64

	packet := make([]byte, 48)
	packet[0] = 0x1B

	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", srv, strconv.Itoa(p)))
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	_, err = conn.Write(packet)
	if err != nil {
		return nil, err
	}

	_, err = conn.Read(packet)
	if err != nil {
		return nil, err
	}

	// retrieve the bytes that we need for the current timestamp
	// data format is unsigned 64 bit long, big endian order
	// see: http://play.golang.org/p/6KRE-2Hq6n
	second = uint64(packet[40])<<24 | uint64(packet[41])<<16 | uint64(packet[42])<<8 | uint64(packet[43])
	fraction = uint64(packet[44])<<24 | uint64(packet[45])<<16 | uint64(packet[46])<<8 | uint64(packet[47])

	nsec := (second * 1e9) + ((fraction * 1e9) >> 32)

	now := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nsec))

	return &now, nil
}
