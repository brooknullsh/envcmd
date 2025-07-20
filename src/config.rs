use std::{
  env::{self},
  fmt::{self, Display},
  fs::{self, File},
  io::Write,
  path::PathBuf,
};

use anyhow::Context;
use serde::{Deserialize, Serialize};

use crate::{abort, log};

#[derive(Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum Kind {
  Directory,
  Branch,
}

impl Display for Kind {
  fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
    match self {
      Kind::Directory => write!(f, "directory"),
      Kind::Branch => write!(f, "branch"),
    }
  }
}

#[derive(Serialize, Deserialize)]
pub struct Config {
  #[serde(rename = "async")]
  pub asynchronous: bool,
  pub kind: Kind,
  pub target: String,
  pub commands: Vec<String>,
}

impl Default for Config {
  fn default() -> Self {
    Self {
      asynchronous: false,
      kind: Kind::Directory,
      target: String::from("foo"),
      commands: Vec::from(["echo 'bar'".into()]),
    }
  }
}

pub fn absolute_config_path() -> PathBuf {
  let Some(path) = env::home_dir().map(|p| p.join(".envcmd/config.json")) else {
    abort!("failed to find the home directory");
  };

  path
}

pub fn read_config_objects(path: &PathBuf) -> anyhow::Result<Vec<Config>> {
  let cfg = fs::read_to_string(path).context("opening configuration")?;
  let cfg: Vec<Config> = serde_json::from_str(&cfg).context("deserialising configuration")?;

  Ok(cfg)
}

pub fn create() -> anyhow::Result<()> {
  let path = absolute_config_path();
  if path.exists() {
    abort!("{} already exists", path.display());
  }

  let Some(parent_dir) = path.parent() else {
    abort!("failed to retrieve parent of {}", path.display());
  };

  fs::create_dir_all(parent_dir).context("creating configuration folder")?;

  let cfg = Vec::from([Config::default()]);
  let cfg = serde_json::to_string_pretty(&cfg).context("serialising default configuration")?;

  File::create(&path)
    .and_then(|mut f| f.write_all(cfg.as_bytes()))
    .context("writing default config")?;

  log!(INFO, "created {}", path.display());
  Ok(())
}

pub fn delete() -> anyhow::Result<()> {
  let path = absolute_config_path();
  if !path.exists() {
    abort!("{} not found", path.display());
  }

  let Some(parent_dir) = path.parent() else {
    abort!("failed to retrieve parent of {}", path.display());
  };

  fs::remove_file(&path).context("removing configuration file")?;
  fs::remove_dir(parent_dir).context("removing configuration folder")?;

  log!(INFO, "deleted {}", path.display());
  Ok(())
}

pub fn list() -> anyhow::Result<()> {
  let path = absolute_config_path();
  if !path.exists() {
    abort!("{} not found", path.display());
  }

  for cfg in read_config_objects(&path)? {
    let status = if cfg.asynchronous { "async" } else { "sync" };
    log!(
      INFO,
      "\x1b[1m{}\x1b[0m ({}) ({status})",
      cfg.target,
      cfg.kind
    );

    cfg.commands.iter().for_each(|c| log!(INFO, "$ {}", c));
  }

  Ok(())
}
