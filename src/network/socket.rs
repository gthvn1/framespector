use crate::network::error::NetworkError;
use std::os::unix::io::RawFd;

pub struct Socket {
    fd: RawFd,
}

impl Socket {
    pub fn new() -> Result<Self, NetworkError> {
        let fd = unsafe { libc::socket(libc::AF_PACKET, libc::SOCK_RAW, 0) };
        if fd < 0 {
            let err = std::io::Error::last_os_error();
            return Err(NetworkError::SocketCreationFailed(err));
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
