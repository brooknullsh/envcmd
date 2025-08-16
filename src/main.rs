use std::{
  env,
  io::{self, BufRead, BufReader},
  process::{Command, Stdio},
  thread,
};

use clap::{Parser, Subcommand};

const DELIMITER_KEY: &str = "ENVCMD_DELIMITER";
const DEFAULT_DELIMITER: &str = ",";

#[macro_export]
macro_rules! log {
  (INFO, $($txt:tt)*) => {
    println!("\x1b[1m\x1b[32mI\x1b[0m {}", format!($($txt)*))
  };
  (WARN, $($txt:tt)*) => {
    println!("\x1b[1m\x1b[33mW\x1b[0m {}", format!($($txt)*))
  };
  (ERROR, $($txt:tt)*) => {
    eprintln!("\x1b[1m\x1b[31mE\x1b[0m {}", format!($($txt)*))
  };
}

#[macro_export]
macro_rules! abort {
  ($($txt:tt)*) => {{
    log!(ERROR, $($txt)*);
    std::process::exit(1)
  }};
}

#[derive(PartialEq)]
pub enum Kind {
  Dir,
  Branch,
}

#[derive(Subcommand, PartialEq)]
#[clap(version)]
enum Commands {
  #[command(about = "View your environment variables", alias = "l")]
  List,
}

#[derive(Parser)]
#[command(about = "Command line tool for running per-environment commands")]
struct Args {
  #[command(subcommand)]
  command: Option<Commands>,
}

/// Normalise a string to a [`Kind`] enum. Since the input string is user-input,
/// default to `None` for misspelt/wrong kind values.
fn normalise_kind(kind: &str) -> Option<Kind> {
  match kind {
    "DIR" => Some(Kind::Dir),
    "BRANCH" => Some(Kind::Branch),
    _ => None,
  }
}

/// Check each [`Kind`]'s target value.
fn kind_matches_target(kind: Kind, target: &str) -> bool {
  (kind == Kind::Dir && dir_match(target)) || (kind == Kind::Branch && branch_match(target))
}

/// Check the [`Kind::Dir`] target value is equal to the current working
/// directory.
fn dir_match(target: &str) -> bool {
  let dir_path = env::current_dir()
    .inspect_err(|err| log!(ERROR, "Reading current directory: {err}"))
    .unwrap();

  dir_path.file_name().is_some_and(|dir| dir == target)
}

/// Check the [`Kind::Branch`] target value is equal to the current working
/// branch.
fn branch_match(target: &str) -> bool {
  let command = Command::new("git")
    .arg("rev-parse")
    .arg("--abbrev-ref")
    .arg("HEAD")
    .stdout(Stdio::piped())
    .spawn()
    .inspect_err(|err| log!(ERROR, "Executing git command: {err}"))
    .unwrap();

  let output = command
    .wait_with_output()
    .inspect_err(|err| log!(ERROR, "Waiting for git command: {err}"))
    .unwrap();

  if output.status.success() {
    return String::from_utf8(output.stdout).is_ok_and(|branch| branch.trim() == target);
  }

  // Quietly handle not being in a git repository. The exit code 128 is a
  // generic fatal error.
  //
  // https://git-scm.com/docs/git-check-ignore.html#Documentation/git-check-ignore.txt-128
  if output.status.code() == Some(128) {
    return false;
  }

  log!(WARN, "Failed to read git branch");
  false
}

/// Iterate through every command from the environment variable value. Spawn a
/// thread for each asynchronous command, block otherwise.
///
/// NOTE: Synchronous commands will spawn two threads for the STDOUT and STDERR
/// streams.
fn run_commands(val: &str, delimiter: &str, asynchronous: bool) {
  let mut handles = Vec::new();

  for (idx, cmd) in val.split(delimiter).enumerate() {
    let cmd = cmd.to_owned();

    if asynchronous {
      let handle = thread::spawn(move || execute_command(&cmd, idx));
      handles.push(handle);
    } else {
      execute_command(&cmd, idx);
    }
  }

  for handle in handles {
    handle.join().unwrap();
  }
}

/// Start the command as a child process, capturing its STDOUT & STDERR and
/// spawning a thread for each.
fn execute_command(cmd: &str, idx: usize) {
  let mut command = Command::new("bash")
    .arg("-c")
    .arg(cmd)
    .stdout(Stdio::piped())
    .stderr(Stdio::piped())
    .spawn()
    .inspect_err(|err| log!(ERROR, "Executing command: {err}"))
    .unwrap();

  let (Some(stdout), Some(stderr)) = (command.stdout.take(), command.stderr.take()) else {
    abort!("Failed to take STDOUT & STDERR");
  };

  thread::scope(|scope| {
    scope.spawn(|| print_stream(stdout, idx));
    scope.spawn(|| print_stream(stderr, idx));
  });
}

/// Reads each line of a stream and prints them. Will rotate through a few
/// colours to distinguish between the commands.
fn print_stream<T>(stream: T, idx: usize)
where
  T: io::Read,
{
  let colours: [&str; 4] = ["\x1b[34m", "\x1b[35m", "\x1b[36m", "\x1b[37m"];

  let colour_index = (idx + 1) % colours.len();
  let colour = colours[colour_index];

  for line in BufReader::new(stream).lines() {
    let Ok(line) = line else {
      log!(WARN, "Failed to read line from stream");
      continue;
    };

    println!("[\x1b[1m{colour}{idx}\x1b[0m] {line}");
  }
}

fn main() {
  let delimiter = env::var(DELIMITER_KEY).unwrap_or(DEFAULT_DELIMITER.to_string());
  let vars = env::vars().filter(|var| var.0.starts_with("ENVCMD") && var.0 != DELIMITER_KEY);

  let args = Args::parse();
  if args.command.is_some() {
    for (key, val) in vars {
      log!(INFO, "\x1b[1m{}\x1b[0m", key);
      val.split(&delimiter).for_each(|cmd| log!(INFO, "$ {cmd}"));
    }

    return;
  }

  for (key, val) in vars {
    let chunks: Vec<&str> = key.split("_").collect();
    let asynchronous = chunks.last().is_some_and(|chunk| *chunk == "ASYNC");

    if asynchronous && chunks.len() != 4 {
      abort!("Invalid format of ENVCMD_<KIND>_<TARGET>_ASYNC: {key}");
    } else if !asynchronous && chunks.len() != 3 {
      abort!("Invalid format of ENVCMD_<KIND>_<TARGET>: {key}");
    }

    // Since we extract the target from a split iterator with an underscore
    // delimiter, target values are limited e.g. directories with dashes in
    // their name. TODO: Better way of parsing target values.
    let (kind, target) = (chunks[1], chunks[2]);
    let Some(kind) = normalise_kind(kind) else {
      abort!("Unrecognised <KIND>: {key}");
    };

    let target = target.to_lowercase();
    if kind_matches_target(kind, &target) {
      run_commands(&val, &delimiter, asynchronous);
    }
  }
}
