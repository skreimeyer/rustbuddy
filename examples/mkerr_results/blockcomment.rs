/*
Welcome to the wild world of block comments
on several lines

This file only has a main function
*/

use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct MyCustomErrorName {
    message: String
}

impl MyCustomErrorName {
    fn new(message: &str) -> MyCustomErrorName {
        MyCustomErrorName{message: message.to_string()}
    }
}

impl fmt::Display for MyCustomErrorName {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f,"{}",self.message)
    }
}

impl Error for MyCustomErrorName {
    fn description(&self) -> &str {
        &self.message
    }
}


use std::f64::consts::PI

fn main() {
    println!("Hello, world!");
    println!("In case you didn't know, pi is {}",PI)
}
