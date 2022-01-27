package analyze

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v3"
	"github.com/dundee/gdu/v5/pkg/fs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	gob.RegisterName("analyze.StoredDir", &StoredDir{})
	gob.RegisterName("analyze.Dir", &Dir{})
	gob.RegisterName("analyze.File", &File{})
}

var DefaultStorage *Storage

type Storage struct {
	db *badger.DB
}

func NewStorage() *Storage {
	st := &Storage{}
	DefaultStorage = st
	return st
}

func (s *Storage) Open() func() {
	options := badger.DefaultOptions("/tmp/badger")
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		panic(err)
	}
	s.db = db

	return func() {
		db.Close()
	}
}

// StoreDir saves item info into badger DB
func (s *Storage) StoreDir(dir fs.Item) error {
	return s.db.Update(func(txn *badger.Txn) error {
		b := &bytes.Buffer{}
		enc := gob.NewEncoder(b)
		err := enc.Encode(dir)
		if err != nil {
			return errors.Wrap(err, "encoding dir value")
		}

		log.Printf("data %s %s", dir.GetPath(), b.String())

		txn.Set([]byte(dir.GetPath()), b.Bytes())
		return nil
	})
}

// LoadDir saves item info into badger DB
func (s *Storage) LoadDir(dir fs.Item) error {
	return s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(dir.GetPath()))
		if err != nil {
			return errors.Wrap(err, "reading stored value")
		}
		return item.Value(func(val []byte) error {
			b := bytes.NewBuffer(val)
			dec := gob.NewDecoder(b)
			return dec.Decode(dir)
		})
	})
}
