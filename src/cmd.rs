use std::{
    env::{self},
    io::{self, BufRead, BufReader},
    process::{Command, Stdio},
    thread::{self},
};

use anyhow::Context;

use crate::{
    abort,
    config::{self, Config, Kind},
    log,
};

pub fn run() -> anyhow::Result<()>
{
    let path = config::absolute_config_path();
    if !path.exists()
    {
        abort!("{} not found", path.display());
    }

    let mut handles = Vec::new();
    let cfg = config::read_config_objects(&path)?
        .into_iter()
        .find(cfg_match)
        .context("no match found")?;

    log!(INFO, "\x1b[1m{}\x1b[0m ({})", cfg.target, cfg.kind);

    for (idx, cmd) in cfg.commands.into_iter().enumerate()
    {
        if cfg.asynchronous
        {
            let handle = thread::spawn(move || execute_command(&cmd, idx));
            handles.push(handle);
            continue;
        }

        execute_command(&cmd, idx);
    }

    for handle in handles
    {
        handle.join().unwrap();
    }

    Ok(())
}

fn cfg_match(cfg: &Config) -> bool
{
    (cfg.kind == Kind::Directory && dir_match(&cfg.target))
        || (cfg.kind == Kind::Branch && branch_match(&cfg.target))
}

fn dir_match(target: &str) -> bool
{
    let dir_path = env::current_dir()
        .inspect_err(|err| log!(ERROR, "reading current directory: {err}"))
        .unwrap();

    dir_path.file_name().is_some_and(|name| name == target)
}

fn branch_match(target: &str) -> bool
{
    let command = Command::new("git")
        .arg("rev-parse")
        .arg("--abbrev-ref")
        .arg("HEAD")
        .stdout(Stdio::piped())
        .spawn()
        .inspect_err(|err| log!(ERROR, "executing git command: {err}"))
        .unwrap();

    let output = command
        .wait_with_output()
        .inspect_err(|err| log!(ERROR, "waiting for git command: {err}"))
        .unwrap();

    if output.status.success()
    {
        return String::from_utf8(output.stdout).is_ok_and(|branch| branch.trim() == target);
    }

    if output.status.code() == Some(128)
    {
        log!(WARN, "no git in current directory");
    }

    log!(WARN, "failed to read git branch");
    false
}

fn execute_command(cmd: &str, idx: usize)
{
    let mut command = Command::new("bash")
        .arg("-c")
        .arg(cmd)
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .inspect_err(|err| log!(ERROR, "executing command: {err}"))
        .unwrap();

    let (Some(stdout), Some(stderr)) = (command.stdout.take(), command.stderr.take())
    else
    {
        abort!("failed to take stdout and stderr");
    };

    thread::scope(|scope| {
        scope.spawn(|| print_stream(stdout, idx));
        scope.spawn(|| print_stream(stderr, idx));
    });
}

fn print_stream(stream: impl io::Read, idx: usize)
{
    let colours: [&str; 4] = ["\x1b[34m", "\x1b[35m", "\x1b[36m", "\x1b[37m"];

    let colour_index = (idx + 1) % colours.len();
    let colour = colours[colour_index];

    for line in BufReader::new(stream).lines()
    {
        let Ok(line) = line
        else
        {
            log!(WARN, "failed to read line from stream");
            continue;
        };

        println!("\x1b[1m{colour}{idx}\x1b[0m {line}");
    }
}
