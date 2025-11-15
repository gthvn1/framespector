Rewrite [network_layers](https://github.com/gthvn1/network_layers) in Go to see how it is
to program in Go.

## Build & Run

- `go build .`
- You need root privileges because program creates a virtual pair socket
  - `sudo ./framespector --help`
- Without parameters the prgram will:
  - Create a **veth0** virtual pair
  - Add **192.168.35.2/24** to **veth0**
  - Listen for incoming frames on **veth0-peer**
- Press `Ctrl-C` to quit, the virtual pair will be cleaned up automatically
- Currently we only print raw frame but you should see ARP if you do `arping -c 192.168.35.3`

```
# sudo ./framespector
time=2025-11-15T17:39:03.434+01:00 level=DEBUG msg="Virtual pair socket created"
time=2025-11-15T17:39:03.438+01:00 level=DEBUG msg="Bind done" iface=veth0-peer
time=2025-11-15T17:39:03.438+01:00 level=INFO msg="Setup network done"
time=2025-11-15T17:39:08.450+01:00 level=INFO msg="frame received" bytes=42
--------- RAW FRAME ---------
ff ff ff ff ff ff 1a a0 bb c8
a5 97 08 06 00 01 08 00 06 04
00 01 1a a0 bb c8 a5 97 c0 a8
23 02 ff ff ff ff ff ff c0 a8
23 03
-----------------------------
^Ctime=2025-11-15T17:39:20.775+01:00 level=INFO msg="ctrl-c received, shutting down..."
time=2025-11-15T17:39:20.866+01:00 level=INFO msg="clean shutdown complete"
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

## POSIX Signals

## Polling (epoll/select)
