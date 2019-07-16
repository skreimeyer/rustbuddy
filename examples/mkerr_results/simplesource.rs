//! Just a comment
//! Second line

use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct simplesourceError {
    message: String
}

impl simplesourceError {
    fn new(message: &str) -> simplesourceError {
        simplesourceError{message: message.to_string()}
    }
}

impl fmt::Display for simplesourceError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f,"{}",self.message)
    }
}

impl Error for simplesourceError {
    fn description(&self) -> &str {
        &self.message
    }
}


fn main() {
    println!("Hello, world!");
    let x = 5;
    let y = 3;
    println!("If you add {} and {} you get {}!", x, y, add(x, y))
}

fn add(a: i32, b: i32) -> i32 {
    return a + b;
}
