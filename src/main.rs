use std::{
  env,
  fmt::{self, Display},
  fs::{File, create_dir_all, remove_file},
  io::{BufRead, BufReader, Write},
  path::PathBuf,
  process,
};

const SEPARATOR: &str = "---";

macro_rules! log {
  (INFO, $($txt:tt)*) => {
    println!("\x1b[1m\x1b[32m[INFO ]\x1b[0m {}", format!($($txt)*))
  };
  (WARN, $($txt:tt)*) => {
    println!("\x1b[1m\x1b[33m[WARN ]\x1b[0m {}", format!($($txt)*))
  };
  (ERROR, $($txt:tt)*) => {
    eprintln!("\x1b[1m\x1b[31m[ERROR]\x1b[0m {}", format!($($txt)*))
  };
  ($($txt:tt)*) => {
    println!("\x1b[1m\x1b[34m[DEBUG]\x1b[0m {}", format!($($txt)*))
  };
}

macro_rules! abort {
  ($($txt:tt)*) => {{
    log!(ERROR, $($txt)*);
    process::exit(1);
  }};
}

macro_rules! ensure {
  ($cond:expr, $($txt:tt)*) => {
    if !$cond {
      abort!($($txt)*);
    }
  }
}

#[derive(Debug)]
struct Config {
  condition: (String, String),
  commands: Vec<String>,
}

impl Default for Config {
  fn default() -> Self {
    Self {
      condition: ("directory".into(), "example".into()),
      commands: vec!["echo 'Hello, world!'".into()],
    }
  }
}

impl Display for Config {
  fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
    write!(
      f,
      "if {} is {}\n{SEPARATOR}\n{}",
      self.condition.0,
      self.condition.1,
      self.commands.join("\n")
    )
  }
}

impl Config {
  fn new() -> Self {
    Self {
      condition: ("".into(), "".into()),
      commands: vec![],
    }
  }
}

fn config_path() -> (PathBuf, PathBuf) {
  let Some(path) = env::home_dir() else {
    abort!("Config: home directory not found");
  };
  (path.join(".envcmd/config.txt"), path.join(".envcmd"))
}

fn create_config() {
  let (file_path, dir_path) = config_path();
  ensure!(!file_path.exists(), "Config: exists at {file_path:?}");

  if let Err(err) = create_dir_all(dir_path)
    .and_then(|_| File::create(&file_path))
    .and_then(|mut file| file.write_all(Config::default().to_string().as_bytes()))
  {
    abort!("Config: creation failed at {file_path:?} due to {err}");
  }
  log!(INFO, "Config: created at {file_path:?}");
}

fn delete_config() {
  let (file_path, _) = config_path();
  ensure!(file_path.exists(), "Config: nothing to delete");

  if let Err(e) = remove_file(&file_path) {
    abort!("Config: failed to delete {file_path:?} due to {e}");
  }
  log!(WARN, "Config: deleted from {file_path:?}");
}

fn view_config() {
  let (file_path, _) = config_path();
  ensure!(file_path.exists(), "Config: nothing to view");

  let Ok(reader) = File::open(&file_path).and_then(|file| Ok(BufReader::new(file))) else {
    abort!("Config: failed to open {file_path:?}");
  };

  let mut lines = reader.lines().peekable();

  let mut is_a_command = false;
  let mut config = Config::new();
  let mut built_configs = Vec::<Config>::new();

  while let Some(line) = lines.next() {
    if let Ok(line) = line {
      if line == SEPARATOR {
        is_a_command = true;
        continue;
      }

      if line.starts_with("if") {
        let parts = line.split_whitespace().collect::<Vec<&str>>();
        config.condition = (parts[1].into(), parts[3].into());
        continue;
      }

      if line.trim().is_empty() || lines.peek().is_none() {
        is_a_command = false;
        if lines.peek().is_none() {
          config.commands.push(line);
        }
        built_configs.push(config);
        config = Config::new();
        continue;
      }

      if is_a_command {
        config.commands.push(line);
      }
    } else {
      abort!("Config: failed to read line from {file_path:?}");
    }
  }

  log!("{built_configs:?}");
}

fn handle_arg(arg: &str) {
  match arg {
    "init" | "-i" => create_config(),
    "delete" | "-d" => delete_config(),
    "view" | "-v" => view_config(),
    _ => abort!("Arguments: unknown value '{arg}'"),
  }
}

fn main() {
  let mut args = env::args();
  args.next(); // Ignore process name

  ensure!(args.len() <= 1, "Arguments: expected 1, got {}", args.len());
  let Some(cmd) = args.next() else {
    abort!("Arguments: not found");
  };

  handle_arg(&cmd);
}
