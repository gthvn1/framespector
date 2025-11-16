Rewrite of [network_layers](https://github.com/gthvn1/network_layers) in Go to
explore how it feels to implement layers in Go.

## Goal

- First goal: reply to ARP request. By default it replies to `arping -c 1 192.168.35.3`
- Next steps: TBD

## Build & Run

- Build the binary: `go build .`
- You must run it as root, because program creates a virtual veth pair:
  - `sudo ./framespector --help`
- Without parameters, the program will:
  - Create a **veth0** virtual ethernet pair
  - Assign **192.168.35.2/24** to **veth0**
  - Listen for incoming frames on **veth0-peer**
- Press `Ctrl-C` to quit, the virtual pair is cleaned up automatically.

```
❯ sudo ./framespector
time=2025-11-16T11:21:09.187+01:00 level=DEBUG msg="proto set" proto=768
time=2025-11-16T11:21:09.191+01:00 level=DEBUG msg="virtual pair socket created"
time=2025-11-16T11:21:09.199+01:00 level=DEBUG msg="bind done" iface=veth0-peer
time=2025-11-16T11:21:09.199+01:00 level=INFO msg="Setup network done"
time=2025-11-16T11:21:09.201+01:00 level=INFO msg="Hit ctrl-c to quit"
time=2025-11-16T11:21:09.201+01:00 level=INFO msg="frame received" bytes=102
time=2025-11-16T11:21:09.202+01:00 level=DEBUG msg="Ethernet: 52:54:00:91:ca:aa -> 52:54:00:38:0e:fe, Type: IPv4, Payload: 88 bytes"
time=2025-11-16T11:21:11.033+01:00 level=INFO msg="frame received" bytes=42
--------- ARP FRAME ---------
ff ff ff ff ff ff 32 f5 c9 0d
d7 15 08 06 00 01 08 00 06 04
00 01 32 f5 c9 0d d7 15 c0 a8
23 02 ff ff ff ff ff ff c0 a8
23 03
-----------------------------
time=2025-11-16T11:21:11.035+01:00 level=DEBUG msg="Ethernet: 32:f5:c9:0d:d7:15 -> ff:ff:ff:ff:ff:ff, Type: ARP, Payload: 28 bytes"
^Ctime=2025-11-16T11:21:14.007+01:00 level=INFO msg="ctrl-c received, shutting down..."
time=2025-11-16T11:21:14.012+01:00 level=INFO msg="stop receiving frame"
time=2025-11-16T11:21:14.012+01:00 level=INFO msg="clean shutdown complete"
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
