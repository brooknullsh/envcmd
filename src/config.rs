use std::{
  fs::{self, File},
  path::PathBuf,
};

use anyhow::{Context, Result};

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
    let config_path = dirs::home_dir()
      .with_context(|| "failed to locate the home directory")?
      .join(&self.path);

    if config_path.exists() {
      log::debug!("config already exists");
      return Ok(());
    }

    let parent_dir = config_path
      .parent()
      .with_context(|| "failed to extract the parent directory")?;

    fs::create_dir_all(parent_dir)?;
    File::create(config_path)?;

    log::debug!("config created");
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

    assert!(!path.exists());
    Config::new(&path).create()?;
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
