package db

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jempe/gopicam/pkg/utils"
)

func newTestDB(t *testing.T) (*DB, func()) {

	// create and remove empty file, we need the file path to create the DB file
	testFile, err := ioutil.TempFile("", "gopicam-dbtest.*.db")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(testFile.Name()) // clean up

	log.Println("create temp db file:", testFile.Name())

	database := &DB{Path: testFile.Name()}

	return database, func() {
		dbPath := database.DBPath()
		database.Close()

		log.Println("removing temp db file:", dbPath)
		os.Remove(dbPath)
	}
}

func TestInitDb(t *testing.T) {
	database, teardown := newTestDB(t)
	defer teardown()

	// Test if DB can be initialized
	t.Run("Init DB", func(t *testing.T) {

		err := database.InitDb()

		if err != nil {
			t.Errorf("error creating database")
		}
	})

	// Test if DB file has been created
	t.Run("DB path exists", func(t *testing.T) {
		dbPath := database.DBPath()

		if !utils.Exists(dbPath) {
			t.Errorf("DB file %s doesn't exist", dbPath)
		}

	})

	t.Run("Non Existent Config Variable", func(t *testing.T) {

		value := database.GetConfigValue("anyvar")

		if value != nil {
			t.Errorf("want nil; got %q", value)
		}
	})

	tests := []struct {
		name string
		key  string
		want []byte
	}{
		{
			name: "Text Variable",
			key:  "text",
			want: []byte("this is the variable value"),
		},
		{
			name: "Numeric Variable",
			key:  "numeric",
			want: []byte("898768768876"),
		},
		{
			name: "Special Chars Variable",
			key:  "special",
			want: []byte("&sdfsdf$@#4234áèsdfü"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := database.SetConfigValue(tt.key, tt.want)

			if err != nil {
				t.Errorf("set config var error %q", tt.want)
			}

			value := database.GetConfigValue(tt.key)

			res := bytes.Compare(value, tt.want)
			if res != 0 {
				t.Errorf("want %q; got %q", tt.want, value)
			}
		})
	}
}
