Rewrite [network_layers](https://github.com/gthvn1/network_layers) in Rust to see how it is
to program in Rust.

## Links

### Raw Sockets
- Official docs:
  - libc crate: https://docs.rs/libc/latest/libc/
  - Search for socket, AF_PACKET, SOCK_RAW
- Better resource:
  - Rust FFI guide: https://doc.rust-lang.org/nomicon/ffi.html
  - Example: https://github.com/rust-lang/libc#usage
- Man pages

## POSIX Signals
- Official docs:
  - libc::sigaction: https://docs.rs/libc/latest/libc/fn.sigaction.html
  - Rust Book chapter on FFI: https://doc.rust-lang.org/book/ch19-01-unsafe-rust.html
- Better approach: Search "rust signal handling" → find signal-hook crate docs or look
at Linux man pages + translate to Rust

## Polling (epoll/select)
- For raw epoll:
  - libc::epoll_create: https://docs.rs/libc/latest/libc/fn.epoll_create.html
  - Linux man pages: man epoll (then translate to Rust)
- For higher-level:
  - mio crate: https://docs.rs/mio/latest/mio/ (async I/O)
  - polling crate: https://docs.rs/polling/latest/polling/ (cross-platform)
