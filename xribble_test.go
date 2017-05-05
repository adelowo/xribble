package xribble

import (
	"reflect"
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

func Test_Xribble_Get(t *testing.T) {
	db := NewXribble()

	i := &Item{"name", []byte("Lanre"), time.Now().Add(time.Minute * 10)}

	if err := db.Add(i); err != nil {
		t.Fatalf(
			`An error occurred while trying to save an item to the database`,
			err)
	}

	newItem, err := db.Get("name")

	if err != nil {
		t.Fatalf(
			`An error occurred while fetching an item from the database`,
			err)
	}

	if !reflect.DeepEqual(i, newItem) {
		t.Fatalf(
			`Items differ.. \n Expected %v \n Got %v`,
			i, newItem)
	}
}

func Test_Xribble_Get_unknownKey(t *testing.T) {
	db := NewXribble()

	_, err := db.Get("unknownKey")

	if err == nil {
		t.Fatal(
			`Error should not contain a nil value 
			because the key does not exist in the database`)
	}

	if err != ErrDatabaseMiss {
		t.Fatalf(
			`Expected error to be a database hit.. Got %v instead`,
			err)
	}
}
