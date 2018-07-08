package idcheck

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestFailure(t *testing.T) {
	assert.True(t, false)
}

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

func TestService_NewID_NewID(t *testing.T) {
	t.Run("hash byte matched known good value", func(t *testing.T) {
		bts := make([]byte, 15)
		rdr := bytes.NewReader(bts)
		svc := NewService(WithReader(rdr))
		id, err := svc.NewID()
		assert.Nil(t, err)
		assert.Equal(t, bts, []byte(id)[:15])
		assert.Equal(t, uint8(0xec), []byte(id)[15])
	})
}

func TestFromBase64(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		id, err := FromBase64("")
		assert.Nil(t, err)
		assert.Equal(t, make(ID, 0), id)
	})

	t.Run("straight a", func(t *testing.T) {
		id, err := FromBase64("AAAAAAAAAAAAAAAAAAAAAA")
		assert.Nil(t, err)
		assert.Equal(t, make(ID, 16), id)
	})
}

func TestID_Base64(t *testing.T) {
	assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAA", make(ID, 16).Base64())
	id, err := NewID()
	assert.Nil(t, err)
	idStr := id.Base64()
	id2, err := FromBase64(idStr)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
}

func TestService_ValidID(t *testing.T) {
	t.Run("returns true for a valid ID", func(t *testing.T) {
		idBytes := make([]byte, 16)
		idBytes[15] = 0xec
		id := ID(idBytes)
		assert.True(t, ValidID(id))
	})

	t.Run("returns false when hash doesn't match", func(t *testing.T) {
		idBytes := make([]byte, 16)
		idBytes[15] = 0x01
		id := ID(idBytes)
		assert.False(t, ValidID(id))
	})

	t.Run("returns false for empty id", func(t *testing.T) {
		var id ID
		assert.False(t, ValidID(id))
	})

	t.Run("returns false when id is too long", func(t *testing.T) {
		idBytes := make([]byte, 16)
		idBytes[15] = 0xec
		id := ID(idBytes)
		assert.False(t, NewService(WithIDLength(15)).ValidID(id))
	})

	t.Run("returns false when id is too short", func(t *testing.T) {
		idBytes := make([]byte, 16)
		idBytes[15] = 0xec
		id := ID(idBytes)
		assert.False(t, NewService(WithIDLength(17)).ValidID(id))
	})
}
