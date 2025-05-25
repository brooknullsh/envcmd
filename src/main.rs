use std::{
  env,
  fs::{File, create_dir_all},
  path::PathBuf,
  process,
};

macro_rules! log {
  (INFO, $($message:tt)*) => {
    println!("\x1b[1m\x1b[32m[INFO ]\x1b[0m {}", format!($($message)*))
  };
  (WARN, $($message:tt)*) => {
    println!("\x1b[1m\x1b[33m[WARN ]\x1b[0m {}", format!($($message)*))
  };
  (ERROR, $($message:tt)*) => {
    eprintln!("\x1b[1m\x1b[31m[ERROR]\x1b[0m {}", format!($($message)*))
  };
  ($($message:tt)*) => {
    println!("\x1b[1m\x1b[34m[DEBUG]\x1b[0m {}", format!($($message)*))
  };
}

macro_rules! abort {
  ($($message:tt)*) => {{
    log!(ERROR, $($message)*);
    process::exit(1);
  }};
}

macro_rules! ensure {
  ($condition:expr, $($message:tt)*) => {
    if !$condition {
      abort!($($message)*);
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
  let (path, dir_path) = config_path();
  ensure!(!path.exists(), "Config: exists at {}", path.display());

  log!("Config: creating at {}", path.display());
  if let Err(err) = create_dir_all(dir_path).and_then(|_| File::create(&path)) {
    abort!("Config: creation failed, {}", err);
  };
  log!(INFO, "Config: created at {}", path.display());
}

fn handle_arg(arg: String) {
  match arg.as_str() {
    "init" | "-i" => create_config(),
    _ => abort!("Arguments: unknown value '{}'", arg),
  }
}

fn main() {
  let mut args = env::args();
  args.next(); // Ignore process name

  ensure!(args.len() <= 1, "Arguments: expected 1, got {}", args.len());
  let Some(cmd) = args.next() else {
    return log!(WARN, "Arguments: not found");
  };

  handle_arg(cmd);
}
