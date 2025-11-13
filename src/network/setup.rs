use crate::network::error;
use std::process::Command;

fn run_command(cmd: &str, args: &[&str]) -> Result<(), error::NetworkError> {
    let output =
        Command::new(cmd)
            .args(args)
            .output()
            .map_err(|_| error::NetworkError::CommandFailed {
                cmd: cmd.to_string(),
                msg: None,
            })?;

    if !output.status.success() {
        let msg = format!("{}", String::from_utf8_lossy(&output.stderr),);
        return Err(error::NetworkError::CommandFailed {
            cmd: cmd.to_string(),
            msg: Some(msg),
        });
    }

    Ok(())
}

// This create virtual pair
pub fn create_veth() -> Result<(), error::NetworkError> {
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
