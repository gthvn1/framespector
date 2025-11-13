use std::fmt;

pub enum NetworkError {
    CommandFailed { cmd: String, msg: Option<String> },
    SocketCreationFailed(std::io::Error),
}

impl fmt::Display for NetworkError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            NetworkError::CommandFailed { cmd, msg } => {
                if let Some(m) = msg {
                    write!(f, "command '{cmd}' failed: {m}")
                } else {
                    write!(f, "command '{cmd}' failed")
                }
            }
            NetworkError::SocketCreationFailed(e) => write!(f, "failed to create socket: {e}"),
        }
    }
}
