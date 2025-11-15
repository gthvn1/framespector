package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"example.com/framespector/network"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))

	args := ReadArgs()
	if args == nil {
		return
	}

	veth := network.NewVeth(logger, args.Veth, args.IP)

	if err := veth.Setup(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer veth.Cleanup()

	if err := veth.CreateSocket(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := veth.BindPeer(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Setup network done")
	logger.Info("TODO: listen on the socket")
	logger.Info("Waiting 5 seconds before closing...")
	time.Sleep(5 * time.Second)
}

type Args struct {
	Veth string
	IP   string
}

func ReadArgs() *Args {
	// We are expecting --veth <ifacename> and --ip <x.x.x.x/yy>
	// So the virtual pair name and the ip with its subnet
	veth_name := flag.String("veth", "veth0", "Virtual Pair name")
	ip := flag.String("ip", "192.168.35.2/24", "IP address with CIDR")
	help := flag.Bool("help", false, "Print help")

	flag.Parse()

	if *help {
		fmt.Println("Usage: framespector --veth <veth-name> --ip <ip/cidr>")
		flag.PrintDefaults()
		return nil
	}

	if _, _, err := net.ParseCIDR(*ip); err != nil {
		fmt.Printf("%s is not a valid IP address with CIDR\n", *ip)
		return nil
	}
	return &Args{
		Veth: *veth_name,
		IP:   *ip,
	}
}
