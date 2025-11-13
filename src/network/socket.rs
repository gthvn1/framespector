use std::fmt;
use std::os::unix::io::RawFd;

#[derive(Debug)]
pub enum SocketError {
    SocketCreationFailed(std::io::Error),
}

impl fmt::Display for SocketError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            SocketError::SocketCreationFailed(e) => write!(f, "failed to create socket: {e}"),
        }
    }
}

impl std::error::Error for SocketError {}

pub struct Socket {
    fd: RawFd,
}

impl Socket {
    pub fn new() -> Result<Self, SocketError> {
        let fd = unsafe { libc::socket(libc::AF_PACKET, libc::SOCK_RAW, 0) };
        if fd < 0 {
            let err = std::io::Error::last_os_error();
            return Err(SocketError::SocketCreationFailed(err));
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
