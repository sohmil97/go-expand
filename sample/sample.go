package sample

import (
	"x/dsl"
)

func main() {
	dsl.Log(23, `SELECT * FROM users WHERE id = ?`, "a")
}
