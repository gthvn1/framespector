use std::os::unix::io::RawFd;

fn main() {
    let _sockfd = match create_socket() {
        Ok(sockfd) => sockfd,
        Err(e) => {
            println!("{e}");
            return;
        }
    };

    println!("Socket created");
}

fn create_socket() -> Result<i32, String> {
    // https://docs.rs/libc/latest/libc/
    // AF_PACKET -> Low-level packet interface (man 7 packet)
    //           == communication domain; this is the protocol family
    // SOCK_RAW -> raw network protocol access
    let sockfd: RawFd = unsafe { libc::socket(libc::AF_PACKET, libc::SOCK_RAW, 0) };
    if sockfd == -1 {
        return Err(String::from("Failed to create the socket"));
    }

    Ok(sockfd)
}
