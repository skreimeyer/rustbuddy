//! Just a comment
//! Second line

fn main() {
    println!("Hello, world!");
    let x = 5;
    let y = 3;
    println!("If you add {} and {} you get {}!", x, y, add(x, y))
}

fn add(a: i32, b: i32) -> i32 {
    return a + b;
}
