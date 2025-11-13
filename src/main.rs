use framespector::network::setup;
use framespector::network::socket::Socket;

fn main() {
    if let Err(e) = setup_network() {
        eprintln!("{e}");
        std::process::exit(1);
    }

    println!("Network setup done");
}

fn setup_network() -> Result<Socket, Box<dyn std::error::Error>> {
    setup::create_veth()?;

    let sockfd = Socket::new()?;
    println!("Socket created");

    Ok(sockfd)
}
