package tools

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

// KeyValueStore provides an in-memory key/value store which is persisted to
// a file. The file handle itself is not kept locked for the duration; it is
// only locked during load and save, to make it concurrency friendly. When
// saving, the store uses optimistic locking to determine whether the db on disk
// has been modified by another process; in which case it loads the latest
// version and re-applies modifications made during this session. This means
// the Lost Update db concurrency issue is possible; so don't use this if you
// need more DB integrity than Read Committed isolation levels.
type KeyValueStore struct {
	mu       sync.RWMutex
	filename string
	log      []keyValueChange

	// This is the persistent data
	// version for optimistic locking, this field is incremented with every Save()
	version int64
	db      map[string]interface{}
}

type keyValueOperation int

const (
	// Set a value for a key
	keyValueSetOperation = keyValueOperation(iota)
	// Removed a value for a key
	keyValueRemoveOperation = keyValueOperation(iota)
)

type keyValueChange struct {
	operation keyValueOperation
	key       string
	value     interface{}
}

// NewKeyValueStore creates a new store and initialises it with contents from
// the named file, if it exists
func NewKeyValueStore(filepath string) (*KeyValueStore, error) {
	kv := &KeyValueStore{filename: filepath, db: make(map[string]interface{})}
	err := kv.loadAndMergeIfNeeded()
	return kv, err
}

// Set updates the key/value store in memory
// Changes are not persisted until you call Save()
func (k *KeyValueStore) Set(key string, value interface{}) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.db[key] = value
	k.log = append(k.log, keyValueChange{keyValueSetOperation, key, value})
}

// Remove removes the key and its value from the store in memory
// Changes are not persisted until you call Save()
func (k *KeyValueStore) Remove(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	delete(k.db, key)
	k.log = append(k.log, keyValueChange{keyValueRemoveOperation, key, nil})
}

// Get retrieves a value from the store, or nil if it is not present
func (k *KeyValueStore) Get(key string) interface{} {
	// Read-only lock
	k.mu.RLock()
	defer k.mu.RUnlock()

	// zero value of interface{} is nil so this does what we want
	return k.db[key]
}

// Save persists the changes made to disk
// If any changes have been written by other code they will be merged
func (k *KeyValueStore) Save() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Short-circuit if we have no changes
	if len(k.log) == 0 {
		return nil
	}

	// firstly peek at version; open read/write to keep lock between check & write
	f, err := os.OpenFile(k.filename, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	stat, _ := os.Stat(k.filename)

	defer f.Close()

	// Only try to merge if > 0 bytes, ignore empty files (decoder will fail)
	if stat.Size() > 0 {
		k.loadAndMergeReaderIfNeeded(f)
		// Now we overwrite the file
		f.Seek(0, os.SEEK_SET)
	}

	k.version++

	enc := gob.NewEncoder(f)
	err = enc.Encode(k.version)
	if err != nil {
		return fmt.Errorf("Error while writing version data to %v: %v", k.filename, err)
	}
	err = enc.Encode(k.db)
	if err != nil {
		return fmt.Errorf("Error while writing new key/value data to %v: %v", k.filename, err)
	}
	// Clear log now that it's saved
	k.log = nil

	return nil
}

// Reads as little as possible from the passed in file to determine if the
// contents are different from the version already held. If so, reads the
// contents and merges with any outstanding changes. If not, stops early without
// reading the rest of the file
func (k *KeyValueStore) loadAndMergeIfNeeded() error {
	stat, err := os.Stat(k.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // missing is OK
		}
		return err
	}
	// Do nothing if empty file
	if stat.Size() == 0 {
		return nil
	}

	f, err := os.OpenFile(k.filename, os.O_RDONLY, 0664)
	if err == nil {
		defer f.Close()
		return k.loadAndMergeReaderIfNeeded(f)
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// As loadAndMergeIfNeeded but lets caller decide how to manage file handles
func (k *KeyValueStore) loadAndMergeReaderIfNeeded(f io.Reader) error {
	var versionOnDisk int64
	// Decode *only* the version field to check whether anyone else has
	// modified the db; gob serializes structs in order so it will always be 1st
	dec := gob.NewDecoder(f)
	err := dec.Decode(&versionOnDisk)
	if err != nil {
		return fmt.Errorf("Problem checking version of key/value data from %v: %v", k.filename, err)
	}
	// Totally uninitialised Version == 0, saved versions are always >=1
	if versionOnDisk != k.version {
		// Reload data & merge
		var dbOnDisk map[string]interface{}
		err = dec.Decode(&dbOnDisk)
		if err != nil {
			return fmt.Errorf("Problem reading updated key/value data from %v: %v", k.filename, err)
		}
		k.reapplyChanges(dbOnDisk)
	}
	return nil
}

// reapplyChanges replays the changes made since the last load onto baseDb
// and stores the result as our own DB
func (k *KeyValueStore) reapplyChanges(baseDb map[string]interface{}) {
	for _, change := range k.log {
		switch change.operation {
		case keyValueSetOperation:
			baseDb[change.key] = change.value
		case keyValueRemoveOperation:
			delete(baseDb, change.key)
		}
	}
	// Note, log is not cleared here, that only happens on Save since it's a
	// list of unsaved changes
	k.db = baseDb

}

// RegisterTypeForKeyValueStorage registers a custom type (e.g. a struct) for
// use in the key value store. This is necessary if you intend to pass custom
// structs to KeyValueStore.Set() rather than primitive types.
func RegisterTypeForKeyValueStorage(val interface{}) {
	gob.Register(val)
}