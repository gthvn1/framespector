Rewrite of [network_layers](https://github.com/gthvn1/network_layers) in Go to
explore how it feels to implement layers in Go.

## Goals

- [x] reply to ARP request. By default it replies to `arping -c 1 192.168.35.3`
- [x] parse IPv4 packet
- [ ] handle ICMP protocol
- Next steps: TBD

## Build & Run

- Build the binary: `go build .`
- You must run it as root, because program creates a virtual veth pair:
  - `sudo ./framespector --help`
- Without parameters, the program will:
  - Create a **veth0** virtual ethernet pair
  - Assign **192.168.35.2/24** to **veth0**
  - Listen for incoming frames on **veth0-peer**
    - By default peer responds to arping **192.168.35.3**
- Press `Ctrl-C` to quit, the virtual pair is cleaned up automatically.

```
❯ sudo ./framespector
time=2025-11-18T13:12:10.711+01:00 level=DEBUG msg="proto set" proto=768
time=2025-11-18T13:12:10.711+01:00 level=DEBUG msg="virtual pair socket created"
time=2025-11-18T13:12:10.717+01:00 level=DEBUG msg="bind done" iface=veth0-peer
time=2025-11-18T13:12:10.717+01:00 level=INFO msg="Setup network done"
time=2025-11-18T13:12:10.717+01:00 level=INFO msg="Hit ctrl-c to quit"
time=2025-11-18T13:12:10.753+01:00 level=INFO msg="frame received" bytes=90
time=2025-11-18T13:12:10.753+01:00 level=WARN msg=todo what="handle IPv6 frame" type="IPv6 (0x86DD)"
time=2025-11-18T13:12:13.522+01:00 level=INFO msg="frame received" bytes=42
--------- ARP FRAME ---------
ff ff ff ff ff ff 72 f5 a3 66
05 1e 08 06 00 01 08 00 06 04
00 01 72 f5 a3 66 05 1e c0 a8
23 00 ff ff ff ff ff ff c0 a8
23 03
-----------------------------
^Ctime=2025-11-18T13:12:16.962+01:00 level=INFO msg="ctrl-c received, shutting down..."
time=2025-11-18T13:12:16.974+01:00 level=INFO msg="stop receiving frame"
time=2025-11-18T13:12:16.975+01:00 level=INFO msg="clean shutdown complete"
```

## Tools we are using
- [Download GO](https://go.dev/dl/)
- [GoPLS](https://go.dev/gopls/)
- [Staticcheck](https://staticcheck.dev/)
- [Goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [Delve](https://github.com/go-delve/delve)
- We use [air](https://github.com/air-verse/air) for quick feedback

## Links

- https://pkg.go.dev/std
- https://gobyexample.com/
- https://go.dev/doc/
- https://go.dev/doc/tutorial/handle-errors

### Raw Sockets

- To access raw syscall we use [sys/unix](https://pkg.go.dev/golang.org/x/sys/unix)

