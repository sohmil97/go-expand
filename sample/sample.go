package sample

import (
	"x/dsl"
)

func main() {
	dsl.Log(5, `SELECT * FROM users WHERE id = ?`, 1)
}
