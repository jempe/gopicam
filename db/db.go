package db

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

const DB_VERSION = 1
const logTag = "BoltDB:"

type DB struct {
	Db   *bolt.DB
	Path string
}

type Filters struct {
	Operator   string
	Conditions []Condition
}

type Condition struct {
	Field      string
	Comparison string
	Value      interface{}
}

type SortBy struct {
	Field     string
	Direction string
}

func (boltdb *DB) InitDb() error {
	var err error
	dbPath := boltdb.Path

	boltdb.Db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Println(logTag, "error creating DB")
		log.Println(err)
		return err
	}

	if boltdb.Db != nil {
		buckets := []string{"devices", "locations", "photos", "videos", "audios", "requests", "configuration"}

		for _, bucket := range buckets {
			err = boltdb.createBucket(bucket)
			if err != nil {
				log.Println(logTag, "error creating", bucket, "bucket")
				log.Println(err)
				return err
			}
		}
	}

	var version int

	migrationsValue := boltdb.GetConfigValue("migrations")

	if migrationsValue != nil {
		version, err = strconv.Atoi(string(migrationsValue))

		if err != nil {
			log.Println(logTag, "error reading current migration")
			log.Println(err)
			return err

		}
	}

	if DB_VERSION != version {
		err = boltdb.SetConfigValue("migrations", []byte(strconv.Itoa(DB_VERSION)))

		if err != nil {
			log.Println(logTag, "error setting current migration")
			log.Println(err)
			return err
		}
	}

	return err
}

// createBucket creates a new bucket in the boltdb Database only if doesn't exist
//
func (boltdb *DB) createBucket(bucketName string) error {
	err := boltdb.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return err
}

func (boltdb *DB) GetConfigValue(variable string) []byte {
	var value []byte

	_ = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("configuration"))
		v := b.Get([]byte(variable))

		value = v

		return nil
	})

	return value
}

func (boltdb *DB) SetConfigValue(variable string, value []byte) error {
	err := boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("configuration"))
		err := b.Put([]byte(variable), value)
		return err
	})

	return err
}

func andQuery(query string, totalArguments int) string {
	if totalArguments > 0 {
		return query + " AND"
	} else {
		return query
	}
}

func orQuery(query string, totalArguments int) string {
	if totalArguments > 0 {
		return query + " OR"
	} else {
		return query
	}
}

func emptyOrContains(haystack []string, needle string) bool {
	if len(haystack) == 0 {
		return true
	}

	for _, field := range haystack {
		if field == needle {
			return true
		}
	}
	return false
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi converts byte to int
//
func btoi(v []byte) int {
	i := int(binary.BigEndian.Uint64(v[:]))
	return i
}

func (boltdb *DB) Close() {
	boltdb.Db.Close()
}

func (boltdb *DB) DBPath() string {
	return boltdb.Db.Path()
}
