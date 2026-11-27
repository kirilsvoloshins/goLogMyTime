package utils

func ClampCursor(cursor, length int) int {
	if length == 0 || cursor < 0 {
		return 0
	}

	if cursor >= length {
		return length - 1
	}

	return cursor
}
