/// Garbage

enum FlashMessage {
	Success, // A unit variant
	Warning{ category: i32, message: String }, // A struct variant
	Error(String) // A tuple variant
  }

also garbage.
