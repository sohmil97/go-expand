package sample

import (
	"x/dsl"
)

func main() {
	dsl.Log(`SELECT * FROM users WHERE id = ?`, 1)
}
