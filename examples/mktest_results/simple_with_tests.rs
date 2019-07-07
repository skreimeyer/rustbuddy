fn main() {
    println!("Hello, world!");
    let x = 5;
    let y = 3;
    println!("If you add {} and {} you get {}!", x, y, add(x, y))
}

fn add(a: i32, b: i32) -> i32 {
    return a + b;
}

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	fn test_add() {
		struct Input {
			a: i32,
			 b: i32,
			};
		struct Output {
			R: i32
		};
		struct Case { 
			inpt: Input,
			out: Output,
			comment: string,
		};
		// __TEST CASES GO HERE__
		let cases = vec![
			// FIXME
			// Case { 
			// 	inpt: Input {},
			// 	out: Output{},
			// 	comment: "",
			// },
		]
		// __END TEST CASES__
		for c in cases.iter() {
			assert!(add(c.inpt.a,c.inpt.b,) == c.out, c.comment)
		}
	}
	
}
