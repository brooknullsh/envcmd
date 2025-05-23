use serde::{Deserialize, Serialize};
use serde_json::{from_str, to_string_pretty};
use std::{
  fmt::{self, Display, Formatter},
  fs::{File, create_dir_all, read_to_string, remove_dir, remove_file},
  io::Write,
  path::PathBuf,
};

use anyhow::{Result, bail, ensure};

#[derive(Serialize, Deserialize)]
struct Content {
  condition: Vec<String>,
  commands: Vec<String>,
}

impl Default for Content {
  fn default() -> Self {
    Self {
      condition: vec!["directory".into(), "example".into()],
      commands: vec!["echo 'Hello, world!'".into()],
    }
  }
}

impl Display for Content {
  fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
    write!(
      f,
      "\nif {}\n---\n{}\n",
      self.condition.join(" is "),
      self.commands.join("\n")
    )
  }
}

pub struct Config {
  whole_path: PathBuf,
  dir_path: PathBuf,
}

impl Config {
  pub fn new() -> Result<Self> {
    let Some(home_path) = dirs::home_dir() else {
      bail!("failed to find your home directory");
    };

    Ok(Self {
      whole_path: home_path.join(".envcmd/config.json"),
      dir_path: home_path.join(".envcmd"),
    })
  }

  pub fn create(&mut self) -> Result<()> {
    ensure!(!self.whole_path.exists(), "configuration already exists");

    create_dir_all(&self.dir_path)?;
    let mut file = File::create(&self.whole_path)?;

    let default_content = Content::default();
    let default_json = to_string_pretty(&default_content)?;
    file.write_all(default_json.as_bytes())?;

    log::info!(
      "config created at {}{}",
      self.whole_path.display(),
      default_content
    );
    Ok(())
  }

  pub fn delete(&mut self) -> Result<()> {
    ensure!(self.whole_path.exists(), "no configuration found");

    remove_file(&self.whole_path)?;
    remove_dir(&self.dir_path)?;

    log::info!("config deleted from {}", self.whole_path.display());
    Ok(())
  }

  pub fn view(&mut self) -> Result<()> {
    ensure!(self.whole_path.exists(), "no configuration found");

    let content_str = read_to_string(&self.whole_path)?;
    let content: Content = from_str(&content_str)?;

    log::info!("{}", content);
    Ok(())
  }
}
