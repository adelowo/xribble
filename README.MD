### Xribble - A tiny XML database in Golang 


Just a troll project to making myself familiar with the XML package in Golang

- Items can be created
  - Supports Encryption
- Items can be deleted
- The database can be dropped

````go

package main

import (
	"fmt"
	"time"
	
	"github.com/adelowo/xribble"
)

func main() {

	x := xribble.NewXribble(xribble.Encrypt(true))

	fmt.Println(
		x.Add(&xribble.Item{"name", []byte("Lanre"),
			time.Now().Add(time.Minute * 10)}))

	fmt.Println(x.Get("name"))
	fmt.Println(x.Delete("name"))
	
	//Add some more data
	fmt.Print(x.Drop()) //Drop the database
}
````


Other available options for the Database factory (`NewXribble`) are

- `Basedir` : Set the base directory of the database
 
- `FS` : Set up a custom implementation for io operations

- `Encrypt` : Should data be encrypted or not.

> Requires you have `XRIBBLE_KEY` in your environment.. And it must be a byte array (with it's length as 32). This is simply because encryption is based on `AES-256`.

