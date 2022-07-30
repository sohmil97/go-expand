package sample

import (
	"x/dsl"
)

func test() {
	dsl.Query(`SELECT * FROM users WHERE id = ?`, 1)
}
