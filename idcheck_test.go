package idcheck

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"testing"
)

func TestHashTable(t *testing.T) {
	t.Run("contains all uint8 values", func(t *testing.T) {
		for i := 0; i <= math.MaxUint8; i++ {
			found := false
			for _, n := range hashTable {
				if n == uint8(i) {
					found = true
					break
				}
			}
			if !found {
				t.Logf("missing value: %d", i)
				t.Fail()
			}
		}
	})

	t.Run("is sufficiently shuffled", func(t *testing.T) {
		for i, v := range hashTable {
			if uint8(i) == v {
				t.Logf("table position and value coincide at %d", i)
			}
		}
	})
}

func assertNil(t *testing.T, v interface{}) {
	t.Helper()
	assertEqual(t, nil, v)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Logf("values are not equal\nexpected: %v\ngot: %v", expected, actual)
		t.Fail()
	}
}

func assert(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Log("expected true but got false")
		t.Fail()
	}
}

func TestNewID(t *testing.T) {
	t.Run("hash byte matched known good value", func(t *testing.T) {
		bts := make([]byte, hashByte)
		idChecker := NewIDChecker(Reader(bytes.NewReader(bts)), Salt("pepper"))
		want := &ID{}
		want[hashByte] = 208
		id, err := idChecker.NewID()
		assertNil(t, err)
		assertEqual(t, want, id)
	})

	t.Run("hash byte changes with the salt", func(t *testing.T) {
		bts := make([]byte, 30)
		idChecker := newIDChecker(Reader(bytes.NewReader(bts)), Salt("salt"))
		id, err := idChecker.NewID()
		assertNil(t, err)
		idChecker.SetSalt("pepper")
		id2, err := idChecker.NewID()
		assertNil(t, err)
		assert(t, !bytes.Equal(id[:], id2[:]))
	})

	t.Run("errors on error reader", func(t *testing.T) {
		idChecker := newIDChecker(Reader(&errReader{}))
		_, err := idChecker.NewID()
		assertEqual(t, errExample, err)
	})
}

func TestFromBase64(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		_, err := FromBase64("")
		assertEqual(t, "str is not the correct length", err.Error())
	})

	t.Run("straight a", func(t *testing.T) {
		id, err := FromBase64("AAAAAAAAAAAAAAAAAAAAAA")
		assertNil(t, err)
		assertEqual(t, &ID{}, id)
	})

	t.Run("invalid string", func(t *testing.T) {
		_, err := FromBase64("\\\\")
		assertEqual(t, "illegal base64 data at input byte 0", err.Error())
	})
}

func TestID_Base64(t *testing.T) {
	assertEqual(t, "AAAAAAAAAAAAAAAAAAAAAA", (&ID{}).Base64())
	id, err := NewID()
	assertNil(t, err)
	idStr := id.Base64()
	id2, err := FromBase64(idStr)
	assertNil(t, err)
	assertEqual(t, id, id2)
}

func TestValidID(t *testing.T) {
	t.Run("returns true for a valid ID", func(t *testing.T) {
		var id ID
		id[hashByte] = 0xec
		assert(t, ValidID(&id))
	})

	t.Run("returns false when hash doesn't match", func(t *testing.T) {
		var id ID
		id[hashByte] = 0xe1
		assert(t, !ValidID(&id))
	})

	t.Run("returns false for empty id", func(t *testing.T) {
		var id ID
		assert(t, !ValidID(&id))
	})
}

var errExample = errors.New("an error")

type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errExample
}
