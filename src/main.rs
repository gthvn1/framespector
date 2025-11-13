use framespector::network::setup;
use framespector::network::socket::Socket;

fn main() {
    if let Err(e) = setup_network() {
        eprintln!("{e}");
        std::process::exit(1);
    }

    let _sockfd = match Socket::new() {
        Ok(sockfd) => sockfd,
        Err(e) => {
            println!("{e}");
            return;
        }
    };

    println!("Socket created");
}

fn setup_network() -> Result<(), Box<dyn std::error::Error>> {
    setup::create_veth()?;
    Ok(())
}
