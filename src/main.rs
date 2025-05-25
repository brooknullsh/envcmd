use std::{
  env,
  fs::{File, create_dir_all},
  path::PathBuf,
  process,
};

macro_rules! log {
  (INFO, $($message:tt)*) => {
    println!("\x1b[1m\x1b[32m[INFO ]\x1b[0m {}", format!($($message)*));
  };
  (WARN, $($message:tt)*) => {
    println!("\x1b[1m\x1b[33m[WARN ]\x1b[0m {}", format!($($message)*));
  };
  (ERROR, $($message:tt)*) => {
    eprintln!("\x1b[1m\x1b[31m[ERROR]\x1b[0m {}", format!($($message)*));
  };
  ($($message:tt)*) => {
    println!("\x1b[1m\x1b[34m[DEBUG]\x1b[0m {}", format!($($message)*));
  };
}

macro_rules! abort {
  ($($message:tt)*) => {
    log!(ERROR, $($message)*);
    process::exit(1);
  };
}

macro_rules! ensure {
  ($condition:expr, $($message:tt)*) => {
    if !$condition {
      abort!($($message)*);
    }
  }
}

fn config_path() -> PathBuf {
  let Some(config_path) = env::home_dir().map(|path| path.join(".envcmd/config.txt")) else {
    abort!("Config: home directory not found");
  };
  config_path
}

fn create_config() {
  let config_path = config_path();
  let Some(parent_dir) = config_path.parent() else {
    abort!("Config: parent directory not found");
  };

  if let Err(err) = create_dir_all(parent_dir).and_then(|_| File::create(&config_path)) {
    abort!(
      "Config: failed to create config at {}, {}",
      config_path.display(),
      err
    );
  }

  log!(INFO, "Config: created at {}", config_path.display());
}

fn handle_arg(arg: String) {
  match arg.as_str() {
    "init" if !config_path().exists() => create_config(),
    _ => todo!(),
  }
}

fn run() {
  todo!();
}

fn main() {
  let mut args = env::args();
  args.next(); // Ignore process name
  ensure!(args.len() <= 1, "Arguments: expected 1, got {}", args.len());

  let Some(cmd) = args.next() else {
    return run();
  };
  handle_arg(cmd);
}
