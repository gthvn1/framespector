use std::fmt;
use std::process::Command;

#[derive(Debug)]
pub enum SetupError {
    CommandFailed { cmd: String, msg: Option<String> },
}

impl fmt::Display for SetupError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            SetupError::CommandFailed { cmd, msg } => {
                if let Some(m) = msg {
                    write!(f, "command '{cmd}' failed: {m}")
                } else {
                    write!(f, "command '{cmd}' failed")
                }
            }
        }
    }
}

impl std::error::Error for SetupError {}

fn run_command(cmd: &str, args: &[&str]) -> Result<(), SetupError> {
    let output = Command::new(cmd)
        .args(args)
        .output()
        .map_err(|_| SetupError::CommandFailed {
            cmd: cmd.to_string(),
            msg: None,
        })?;

    if !output.status.success() {
        let msg = format!("{}", String::from_utf8_lossy(&output.stderr),);
        return Err(SetupError::CommandFailed {
            cmd: cmd.to_string(),
            msg: Some(msg),
        });
    }

    Ok(())
}

// This create virtual pair
pub fn create_veth() -> Result<(), SetupError> {
    run_command(
        "ip",
        &[
            "link",
            "add",
            "veth0",
            "type",
            "veth",
            "peer",
            "name",
            "veth0-peer",
        ],
    )
}
