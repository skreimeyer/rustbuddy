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

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	fn test_HasAreaRectanglearea() {
		
		struct Output {
			R: i32
		};
		struct Case { 
			obj: Rectangle
			inpt: Input,
			out: Output,
			comment: string,
		};
		// __TEST CASES GO HERE__
		let cases = vec![
			// FIXME
			// Case { 
			//	obj: Rectangle{}
			// 	inpt: Input {},
			// 	out: Output{},
			// 	comment: "",
			// },
		]
		// __END TEST CASES__
		for c in cases.iter() {
			assert!(c.obj.area() == c.out, c.comment)
		}
	}
	#[test]
	fn test_HasAreaCirclearea() {
		
		struct Output {
			R: i32
		};
		struct Case { 
			obj: Circle
			inpt: Input,
			out: Output,
			comment: string,
		};
		// __TEST CASES GO HERE__
		let cases = vec![
			// FIXME
			// Case { 
			//	obj: Circle{}
			// 	inpt: Input {},
			// 	out: Output{},
			// 	comment: "",
			// },
		]
		// __END TEST CASES__
		for c in cases.iter() {
			assert!(c.obj.area() == c.out, c.comment)
		}
	}
	
}
