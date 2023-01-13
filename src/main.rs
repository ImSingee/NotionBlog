use clap::Parser;
use anyhow::Result;

#[derive(Parser)]
struct Args {
    #[clap(short, long)]
    root: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();
    Ok(())
}
