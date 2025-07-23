use clap::{Parser, Subcommand};

mod cmd;
mod config;

#[macro_export]
macro_rules! log
{
  (INFO, $($txt:tt)*) =>
  {
    println!("\x1b[1m\x1b[32mI\x1b[0m {}", format!($($txt)*))
  };
  (WARN, $($txt:tt)*) =>
  {
    println!("\x1b[1m\x1b[33mW\x1b[0m {}", format!($($txt)*))
  };
  (ERROR, $($txt:tt)*) =>
  {
    eprintln!("\x1b[1m\x1b[31mE\x1b[0m {}", format!($($txt)*))
  };
}

#[macro_export]
macro_rules! abort
{
  ($($txt:tt)*) =>
  {{
    log!(ERROR, $($txt)*);
    std::process::exit(1)
  }};
}

#[derive(Subcommand)]
enum Commands
{
  #[command(about = "Create a new configuration", alias = "c")]
  Create,
  #[command(about = "Remove your configuration", alias = "d")]
  Delete,
  #[command(about = "View your configuration", alias = "l")]
  List,
}

#[derive(Parser)]
#[command(about = "Command line tool for running per-environment commands")]
struct Args
{
  #[command(subcommand)]
  command: Option<Commands>,
}

fn main()
{
  let args = Args::parse();

  let out = match args.command
  {
    Some(Commands::Create) => config::create(),
    Some(Commands::Delete) => config::delete(),
    Some(Commands::List) => config::list(),
    None => cmd::run(),
  };

  if let Err(err) = out
  {
    abort!("{:#}", err)
  }
}
