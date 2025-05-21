use anyhow::Result;
use clap::{Parser, Subcommand};
use config::Config;
use env_logger::Builder;
use log::LevelFilter;

mod config;

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

fn main() -> Result<()> {
  init_log();
  let args = Args::parse();

  let config = Config::new(".envcmd/config.json");
  match args.command {
    Commands::Init => config.create()?,
  }

  Ok(())
}
