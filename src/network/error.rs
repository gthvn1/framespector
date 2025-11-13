use std::fmt;

pub enum NetworkError {
    CommandFailed(String),
    CommandFailedWithMsg(String),
    SocketCreationFailed(std::io::Error),
}

impl fmt::Display for NetworkError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            NetworkError::CommandFailed(s) | NetworkError::CommandFailedWithMsg(s) => {
                write!(f, "{s}")
            }
            NetworkError::SocketCreationFailed(e) => write!(f, "failed to create socket: {e}"),
        }
    }
}
