package xribble

import (
	"testing"
	"time"
)

var _ Provider = &XribbleDriver{}

func Test_NewXribblePanicsOnEncounteringAnUnWriteablePath(t *testing.T) {

	defer func() {
		recover()
	}()

	NewXribble(BaseDir("/oops"))
}

func Test_Xribble_Add(t *testing.T) {

	db := NewXribble()

	i := &Item{"name", []byte("Lanre"), time.Now().Add(time.Minute * 10)}

	if err := db.Add(i); err != nil {
		t.Fatalf(
			`An error occurred while trying to save an item to the database`,
			err)
	}
}
