use crate::network::error;
use std::process::Command;

fn run_command(cmd: &str, args: &[&str]) -> Result<(), error::NetworkError> {
    let output = Command::new(cmd)
        .args(args)
        .output()
        .map_err(|_| error::NetworkError::CommandFailed(cmd.to_string()))?;

    if !output.status.success() {
        let msg = format!(
            "stderr: {}\nstdout: {}",
            String::from_utf8_lossy(&output.stderr),
            String::from_utf8_lossy(&output.stdout)
        );
        return Err(error::NetworkError::CommandFailedWithMsg(msg));
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
