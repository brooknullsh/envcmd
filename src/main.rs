use anyhow::Result;
use clap::{Parser, Subcommand};
use config::Config;
use env_logger::Builder;
use log::LevelFilter;

mod config;

#[derive(Subcommand)]
enum Command {
  Init,
  Delete,
  View,
}

#[derive(Parser)]
struct Args {
  #[command(subcommand)]
  command: Command,
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

  let mut config = Config::default();
  match args.command {
    Command::Init => config.create()?,
    Command::Delete => config.delete()?,
    Command::View => config.view()?,
  }

  Ok(())
}
