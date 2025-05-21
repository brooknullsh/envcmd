use std::{
  fs::{self, File},
  path::PathBuf,
};

use anyhow::{Result, anyhow};

pub struct Config {
  path: PathBuf,
}

impl Config {
  pub fn new<T>(path: T) -> Self
  where
    T: Into<PathBuf>,
  {
    Self { path: path.into() }
  }

  pub fn create(&self) -> Result<()> {
    let Some(home_path) = dirs::home_dir() else {
      return Err(anyhow!("Failed to locate the home directory"));
    };

    let config_path = home_path.join(&self.path);
    if config_path.exists() {
      log::debug!("Config already exists");
      return Ok(());
    }

    let Some(parent_dir) = config_path.parent() else {
      return Err(anyhow!("Failed to extract the parent directory"));
    };

    fs::create_dir_all(parent_dir)?;
    File::create(config_path)?;

    log::debug!("Config created");
    Ok(())
  }
}

#[cfg(test)]
mod tests {
  use super::*;

  const PATH: &str = ".test/test.json";

  #[test]
  fn test_create_config_from_scratch() -> Result<()> {
    let path = tempfile::tempdir()?.path().join(PATH);
    let config = Config::new(&path);

    assert!(!path.exists());
    config.create()?;
    assert!(path.exists());

    Ok(())
  }

  #[test]
  fn test_create_config_already_exists() -> Result<()> {
    let path = tempfile::tempdir()?.path().join(PATH);
    let config = Config::new(&path);

    config.create()?;
    assert!(config.create().is_ok());

    Ok(())
  }
}
