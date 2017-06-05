package xribble

import (
	"os"
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
			`An error occurred while trying to save an item to the database...%v`,
			err)
	}
}

func Test_Xribble_Get(t *testing.T) {
	db := NewXribble()

	i := &Item{"name", []byte("Lanre"), time.Now().Add(time.Minute * 10)}

	if err := db.Add(i); err != nil {
		t.Fatalf(
			`An error occurred while trying to save an item to the database..%v`,
			err)
	}

	newItem, err := db.Get("name")

	if err != nil {
		t.Fatalf(
			`An error occurred while fetching an item from the database..%v`,
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

func Test_Xribble_Drop(t *testing.T) {

	db := NewXribble()

	if err := db.Drop(); err != nil {
		t.Fatalf(
			`An error occurred while trying to flush the database..%v`,
			err)
	}

	//Manually assert the directory is gone

	// if ok := db.fs.IsDirectory(db.baseDir); ok {
	// 	t.Fatalf(
	// 		`Flush operation failed .. %v`, ok)
	// 		This should work too
	// }
	if _, err := os.Stat(db.baseDir); os.IsExist(err) {
		t.Fatalf(`Flush operation failed .. %v`, err)
	}

}

func Test_Xribble_Delete(t *testing.T) {
	db := NewXribble()

	db.Add(&Item{"name", []byte("Lanre"), time.Now().Add(time.Hour * 4)})

	if err := db.Delete("name"); err != nil {
		t.Fatalf(
			`An error occurred while trying to delete the key, %s ..%v`,
			"name", err)
	}

	//Try fetching the key

	if _, err := db.Get("name"); err != ErrDatabaseMiss {
		t.Fatalf(
			`Data is supposed not to exist in the database again since
			it was deleted.. %v`, err)
	}
}
