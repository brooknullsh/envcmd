use std::{
  env,
  io::{self, BufRead, BufReader},
  process::{Command, Stdio},
  thread,
};

use clap::{Parser, Subcommand};

const SEPARATOR_KEY: &str = "ENVCMD_SEPARATOR";
const DELIMITER_KEY: &str = "ENVCMD_DELIMITER";
const DEFAULT_SEPARATOR: &str = "-";
const DEFAULT_DELIMITER: &str = ",";

#[macro_export]
macro_rules! log {
  (INFO, $($txt:tt)*) => {
    println!("[\x1b[1m\x1b[32mI\x1b[0m] {}", format!($($txt)*))
  };
  (WARN, $($txt:tt)*) => {
    println!("[\x1b[1m\x1b[33mW\x1b[0m] {}", format!($($txt)*))
  };
  (ERROR, $($txt:tt)*) => {
    eprintln!("[\x1b[1m\x1b[31mE\x1b[0m] {}", format!($($txt)*))
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

/// Only find non-configuration yet related environment variables.
fn include_env_vars(var: &(String, String)) -> bool {
  let exclusions: [&str; 2] = [DELIMITER_KEY, SEPARATOR_KEY];
  let key = var.0.as_str();

  key.starts_with("ENVCMD") && !exclusions.contains(&key)
}

/// Normalise a string to a [`Kind`] enum. Since the input string is user-input,
/// default to `None` for misspelt/wrong [`Kind`] values.
fn normalise_kind(kind: &str) -> Option<Kind> {
  match kind {
    "DIR" => Some(Kind::Dir),
    "BRANCH" => Some(Kind::Branch),
    _ => None,
  }
}

/// Convert the split values of the target to lowercase & join them with the
/// specified separator.
fn normalise_target(target: &[&str], separator: &str) -> String {
  target
    .iter()
    .map(|target| target.to_lowercase())
    .collect::<Vec<String>>()
    .join(separator)
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

/// Iterate through every command from the environment variable value split by
/// the specified delimiter. Spawn a thread for each asynchronous command, block
/// otherwise.
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
  let separator = env::var(SEPARATOR_KEY).unwrap_or(DEFAULT_SEPARATOR.to_string());
  let vars = env::vars().filter(include_env_vars);

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

    if asynchronous && chunks.len() < 4 {
      abort!("Invalid format of ENVCMD_<KIND>_<TARGET>_ASYNC: {key}");
    } else if !asynchronous && chunks.len() < 3 {
      abort!("Invalid format of ENVCMD_<KIND>_<TARGET>: {key}");
    }

    let (kind, target) = (chunks[1], chunks[2]);
    let Some(kind) = normalise_kind(kind) else {
      abort!("Unrecognised <KIND>: {key}");
    };

    // Take all underscore-separated words after the kind until the asynchronous
    // flag or the end to later join them with the set separator char(s). This
    // is to handle targets that aren't simply one word in length, and where
    // said target e.g. a directory name isn't separated by underscores.
    let first_target = &[target];
    let target = if asynchronous {
      chunks.get(2..chunks.len() - 1).unwrap_or(first_target)
    } else {
      chunks.get(2..).unwrap_or(first_target)
    };

    let target = normalise_target(target, &separator);
    if kind_matches_target(kind, &target) {
      run_commands(&val, &delimiter, asynchronous);
    }
  }
}
