package main

func main() {
	const tmpl := `
use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct {{.name}} {
    message: String
}

impl {{.name}} {
    fn new(message: &str) -> {{.name}} {
        {{.name}}{message: message.to_string()}
    }
}

impl fmt::Display for {{.name}} {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f,"{}",self.message)
    }
}

impl Error for {{.name}} {
    fn description(&self) -> &str {
        &self.message
    }
}`
}