use std::env;

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

#[derive(Subcommand, PartialEq)]
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

    let Some(path) = env::home_dir().map(|path| path.join(".envcmd/config.json"))
    else
    {
        abort!("failed to find the home directory");
    };

    if args.command == Some(Commands::Create) && path.exists()
    {
        abort!("{} already exists", path.display());
    }
    else if !path.exists() && args.command != Some(Commands::Create)
    {
        abort!("{} not found", path.display());
    }

    let out = match args.command
    {
        Some(Commands::Create) => config::create(path),
        Some(Commands::Delete) => config::delete(path),
        Some(Commands::List) => config::list(path),
        None => cmd::run(path),
    };

    if let Err(err) = out
    {
        abort!("{:#}", err);
    }
}
