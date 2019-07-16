
use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct SomeErrorName {
    message: String
}

impl SomeErrorName {
    fn new(message: &str) -> SomeErrorName {
        SomeErrorName{message: message.to_string()}
    }
}

impl fmt::Display for SomeErrorName {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f,"{}",self.message)
    }
}

impl Error for SomeErrorName {
    fn description(&self) -> &str {
        &self.message
    }
}

pub struct Point {
    x: i32,
    y: i32,
}

pub struct Rectangle {
    a: Point,
    b: Point,
}

pub struct Circle {
    center: Point,
    radius: i32,
}

pub trait HasArea {
    fn area(&self) -> i32;
}

impl HasArea for Rectangle {
    pub fn area(&self) -> i32 {
        let height = (self.b.y - self.a.y).abs();
        let width = (self.b.x - self.a.x).abs();
        return height * width;
    }
}

impl HasArea for Circle {
    pub fn area(&self) -> i32 {
        let farea = std::f64::consts::PI * f64::from(self.radius).powi(2);
        return std::i32::from(farea);
    }
}
