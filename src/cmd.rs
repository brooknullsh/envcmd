use std::{
  env::{self},
  io::{self, BufRead, BufReader},
  process::{Command, Stdio},
  thread::{self, JoinHandle},
};

use crate::{
  abort,
  config::{self, Kind},
  log,
};

const COLOURS: [&str; 4] = ["\x1b[34m", "\x1b[35m", "\x1b[36m", "\x1b[37m"];

pub fn run() -> anyhow::Result<()> {
  let path = config::absolute_config_path();
  if !path.exists() {
    abort!("{} not found", path.display())
  }

  let mut handles = Vec::new();

  for cfg in config::read_config_objects(&path)? {
    match cfg.kind {
      Kind::Directory if no_match_for_dir(&cfg.target) => continue,
      Kind::Branch if no_match_for_branch(&cfg.target) => continue,
      _ => log!(INFO, "\x1b[1m{}\x1b[0m ({})", cfg.target, cfg.kind),
    }

    process_matched_commands(cfg.commands, cfg.asynchronous, &mut handles);
  }

  for handle in handles {
    handle.join().unwrap();
  }

  Ok(())
}

fn no_match_for_dir(target: &str) -> bool {
  let dir_path = env::current_dir()
    .inspect_err(|e| log!(ERROR, "reading current directory: {e}"))
    .unwrap();

  dir_path.file_name().is_some_and(|n| n != target)
}

fn no_match_for_branch(target: &str) -> bool {
  let command = Command::new("git")
    .arg("rev-parse")
    .arg("--abbrev-ref")
    .arg("HEAD")
    .stdout(Stdio::piped())
    .spawn()
    .inspect_err(|e| log!(ERROR, "executing git command: {e}"))
    .unwrap();

  let output = command
    .wait_with_output()
    .inspect_err(|e| log!(ERROR, "waiting for git command: {e}"))
    .unwrap();

  if output.status.success() {
    return String::from_utf8(output.stdout).is_ok_and(|b| b.trim() != target);
  }

  if output.status.code() == Some(128) {
    log!(WARN, "no git in current directory");
  } else {
    log!(WARN, "failed to read git branch");
  }

  true
}

fn process_matched_commands(
  commands: Vec<String>,
  is_async: bool,
  handles: &mut Vec<JoinHandle<()>>,
) {
  for (idx, cmd) in commands.into_iter().enumerate() {
    if is_async {
      let handle = thread::spawn(move || execute_command(&cmd, idx));
      handles.push(handle);
      continue;
    }

    execute_command(&cmd, idx);
  }
}

fn execute_command(cmd: &str, idx: usize) {
  let mut command = Command::new("bash")
    .arg("-c")
    .arg(cmd)
    .stdout(Stdio::piped())
    .stderr(Stdio::piped())
    .spawn()
    .inspect_err(|e| log!(ERROR, "executing command: {e}"))
    .unwrap();

  let (Some(stdout), Some(stderr)) = (command.stdout.take(), command.stderr.take()) else {
    abort!("failed to take stdout and stderr");
  };

  thread::scope(|s| {
    s.spawn(|| print_stream(stdout, idx));
    s.spawn(|| print_stream(stderr, idx));
  });
}

fn print_stream(stream: impl io::Read, idx: usize) {
  let colour_index = (idx + 1) % COLOURS.len();
  let colour = COLOURS[colour_index];

  for line in BufReader::new(stream).lines() {
    let Ok(line) = line else {
      log!(WARN, "failed to read line from stream");
      continue;
    };

    println!("\x1b[1m{colour}{idx}\x1b[0m {line}");
  }
}
