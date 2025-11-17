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
time=2025-11-16T18:20:09.691+01:00 level=DEBUG msg="proto set" proto=768
time=2025-11-16T18:20:09.691+01:00 level=DEBUG msg="virtual pair socket created"
time=2025-11-16T18:20:09.719+01:00 level=DEBUG msg="bind done" iface=veth0-peer
time=2025-11-16T18:20:09.719+01:00 level=INFO msg="Setup network done"
time=2025-11-16T18:20:09.719+01:00 level=INFO msg="Hit ctrl-c to quit"
time=2025-11-16T18:20:09.719+01:00 level=INFO msg="frame received" bytes=86
time=2025-11-16T18:20:09.719+01:00 level=DEBUG msg="Ethernet: 4c:79:6e:d6:64:1a -> 78:c2:13:1d:e3:50, Type: IPv6, Payload: 72 bytes"
time=2025-11-16T18:20:09.719+01:00 level=DEBUG msg="TODO: decode ipv6"
time=2025-11-16T18:20:12.989+01:00 level=INFO msg="frame received" bytes=58
--------- ARP FRAME ---------
ff ff ff ff ff ff 56 ba 0d b1
f6 af 08 06 00 01 08 00 06 04
00 01 56 ba 0d b1 f6 af c0 a8
23 00 00 00 00 00 00 00 c0 a8
23 03 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
-----------------------------
time=2025-11-16T18:20:12.989+01:00 level=DEBUG msg="Ethernet: 56:ba:0d:b1:f6:af -> ff:ff:ff:ff:ff:ff, Type: ARP, Payload: 44 bytes"
time=2025-11-16T18:20:12.990+01:00 level=INFO msg="sent ARP reply" to_mac=56:ba:0d:b1:f6:af from_mac=ae:a1:97:16:8e:62
time=2025-11-16T18:20:13.045+01:00 level=INFO msg="frame received" bytes=179
time=2025-11-16T18:20:13.045+01:00 level=DEBUG msg="Ethernet: 56:ba:0d:b1:f6:af -> 01:00:5e:00:00:fb, Type: IPv4, Payload: 165 bytes"
time=2025-11-16T18:20:13.045+01:00 level=DEBUG msg="TODO: decode ipv4"
time=2025-11-16T18:20:15.932+01:00 level=INFO msg="frame received" bytes=58
--------- ARP FRAME ---------
ff ff ff ff ff ff 56 ba 0d b1
f6 af 08 06 00 01 08 00 06 04
00 01 56 ba 0d b1 f6 af c0 a8
23 00 00 00 00 00 00 00 c0 a8
23 05 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
-----------------------------
time=2025-11-16T18:20:15.933+01:00 level=DEBUG msg="Ethernet: 56:ba:0d:b1:f6:af -> ff:ff:ff:ff:ff:ff, Type: ARP, Payload: 44 bytes"
time=2025-11-16T18:20:15.933+01:00 level=ERROR msg="ARP request not handled" err="IP 192.168.35.3 is not matching 192.168.35.5"
^Ctime=2025-11-16T18:20:18.774+01:00 level=INFO msg="ctrl-c received, shutting down..."
time=2025-11-16T18:20:18.795+01:00 level=INFO msg="stop receiving frame"
time=2025-11-16T18:20:18.795+01:00 level=INFO msg="clean shutdown complete"
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

