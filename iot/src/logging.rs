use chrono::Local;
use env_logger::{fmt::Color, Builder};
use log::LevelFilter;
use std::io::Write;

pub fn init_logger(log_level: LevelFilter) {
    Builder::new()
        .format(|buf, record| {
            let mut style = buf.style();
            let level_color = match record.level() {
                log::Level::Error => Color::Red,
                log::Level::Warn => Color::Yellow,
                log::Level::Info => Color::Green,
                log::Level::Debug => Color::Blue,
                log::Level::Trace => Color::Cyan,
            };
            style.set_color(level_color);

            writeln!(
                buf,
                "{} [{}] {} - {}{}",
                Local::now().format("%Y-%m-%d %H:%M:%S"),
                style.value(record.level()),
                record.target(),
                record.args(),
                String::new()
            )
        })
        .filter(None, log_level)
        .init();
}

pub fn print_banner() {
    println!(
        r#"
 _______       .-''-.      .-''-.     .-'''-.  _______  .-./`)  ____..--'
\  ____  \   .'_ _   \   .'_ _   \   / _     \\  ____  \\ .-.')|        |
| |    \ |  / ( ` )   ' / ( ` )   ' (`' )/`--'| |    \ |/ `-' \|   .-'  '
| |____/ / . (_ o _)  |. (_ o _)  |(_ o _).   | |____/ / `-'`"`|.-'.'   /
|   _ _ '. |  (_,_)___||  (_,_)___| (_,_). '. |   _ _ '. .---.    /   _/
|  ( ' )  \'  \   .---.'  \   .---..---.  \  :|  ( ' )  \|   |  .'._( )_
| (_{{}};}}_) | \  `-'    / \  `-'    /\    `-'  || (_{{;}}_) ||   |.'  (_'o._)
|  (_,_)  /  \       /   \       /  \       / |  (_,_)  /|   ||    (_,_)|
/_______.'    `'-..-'     `'-..-'    `-...-'  /_______.' '---'|_________|
    "#
    );
}
