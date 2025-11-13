use framespector::network::setup;
use framespector::network::socket::Socket;

fn main() {
    if let Err(e) = setup::create_veth() {
        eprintln!("{e}");
        return;
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
