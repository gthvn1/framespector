use std::os::unix::io::RawFd;

pub struct Socket {
    fd: RawFd,
}

impl Socket {
    pub fn new() -> std::io::Result<Self> {
        let fd = unsafe { libc::socket(libc::AF_PACKET, libc::SOCK_RAW, 0) };
        if fd < 0 {
            return Err(std::io::Error::last_os_error());
        }
        Ok(Socket { fd })
    }

    pub fn fd(&self) -> RawFd {
        self.fd
    }
}

impl Drop for Socket {
    fn drop(&mut self) {
        unsafe {
            libc::close(self.fd);
        }
    }
}

fn main() {
    let _sockfd = match Socket::new() {
        Ok(sockfd) => sockfd,
        Err(e) => {
            println!("{e}");
            return;
        }
    };

    println!("Socket created");
}
