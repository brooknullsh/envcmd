use clap::{Parser, Subcommand};
use env_logger::Builder;
use log::LevelFilter;

#[derive(Subcommand)]
enum Commands {
  Init,
}

#[derive(Parser)]
struct Args {
  #[command(subcommand)]
  command: Commands,
}

fn init_log() {
  Builder::from_default_env()
    .filter_level(LevelFilter::Trace)
    .format_timestamp(None)
    .format_target(false)
    .init();
}

fn main() {
  init_log();
  let args = Args::parse();

  match args.command {
    Commands::Init => log::info!("Hello, world!"),
  }
}
