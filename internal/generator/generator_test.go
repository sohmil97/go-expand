package generator

import (
	"testing"
	"x/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	t.Parallel()

	t.Run("should return empty string with no elements", func(t *testing.T) {
		args := make(types.ArgSpec)
		result := GetArgsString(args)
		assert.Equal(t, "", result)
	})

	t.Run("should return args string with one element", func(t *testing.T) {
		args := make(types.ArgSpec)
		args[types.ParamSpec{
			Identifier: types.Identifier{
				Name: "arg1",
			},
			IsVarArgs: false,
		}] = "first"

		result := GetArgsString(args)
		assert.Equal(t, "first", result)
	})

	t.Run("should return args string with two element", func(t *testing.T) {
		args := make(types.ArgSpec)
		args[types.ParamSpec{
			Identifier: types.Identifier{
				Name: "arg1",
			},
			IsVarArgs: false,
		}] = "first"
		args[types.ParamSpec{
			Identifier: types.Identifier{
				Name: "arg2",
			},
			IsVarArgs: false,
		}] = "second"

		result := GetArgsString(args)
		assert.Equal(t, "first, second", result)
	})

	t.Run("should return args string with varargs", func(t *testing.T) {
		args := make(types.ArgSpec)
		args[types.ParamSpec{
			Identifier: types.Identifier{
				Name: "arg1",
			},
			IsVarArgs: true,
		}] = []interface{}{"first", "second", 3}

		result := GetArgsString(args)
		assert.Equal(t, "first, second, 3", result)
	})

	t.Run("should return args string with one element and varargs", func(t *testing.T) {
		args := make(types.ArgSpec)
		args[types.ParamSpec{
			Identifier: types.Identifier{
				Name: "arg1",
			},
			IsVarArgs: false,
		}] = "first"
		args[types.ParamSpec{
			Identifier: types.Identifier{
				Name: "arg2",
			},
			IsVarArgs: true,
		}] = []interface{}{"second", 3, true}

		result := GetArgsString(args)
		assert.Equal(t, "first, second, 3, true", result)
	})
}
