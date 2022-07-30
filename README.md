# go-expand
Trying a new idea in Go.
What if we could have some form of macro in Go ?
Something that is compatible with current language syntax and does not break IDE
support and will generate high-performance fast code.

## Example
```go
package main

type User struct {
	ID int
	Name string
}

// you write this
func getUser() (*User, error) {
	var user User
	err := db.Query(user, `SELECT * FROM users where ID=?`, 1)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// and it will get expanded to something like this
func getUser() (*User, error) {
	row, err := db.QueryRow(`SELECT * FROM users where ID=?`, 1)
	if err != nil {
		return nil, err
    }
	err = row.Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}
	
	return user, nil
	
}

```

# How it works ?
go-expand is built upon idea of markers, markers are nodes in the AST
of a golang file that are defined by the preprocessor type and they are 
nodes that should be expanded. After finding each one we pass the node to the
processor and it will return a block node that will replace marker node in the AST.
and we then do write the new ast to a file named similar to original file and temporarily rename 
original files to avoid them from being compiled. after the compilation is done
generated files are deleted and we rename original files back.