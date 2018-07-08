package idcheck

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

var hashTable = func() [256]uint8 {
	var vals [256]uint8
	for n, _ := range vals {
		vals[n] = uint8(n)
	}
	for _, i := range []int{2, 3, 5, 7, 11} {
		for n, _ := range vals {
			vals[n], vals[n/i] = vals[n/i], vals[n]
		}
	}
	return vals
}()

func hash(data []byte) (hsh uint8) {
	for _, d := range data {
		hsh = hashTable[hsh^d]
	}
	return
}

type Service struct {
	IDLength int
	idReader io.Reader
}

type ServiceOpt func(*Service)

func WithReader(reader io.Reader) ServiceOpt {
	return func(idSvc *Service) {
		idSvc.idReader = reader
	}
}

func WithIDLength(idLength int) ServiceOpt {
	return func(svc *Service) {
		svc.IDLength = idLength
	}
}

func NewService(opts ...ServiceOpt) *Service {
	svc := &Service{
		idReader: rand.Reader,
		IDLength: 16,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func NewID() (ID, error) {
	return NewService().NewID()
}
func (svc *Service) NewID() (ID, error) {
	idBytes := make([]byte, svc.IDLength-1)
	_, err := svc.idReader.Read(idBytes)
	if err != nil {
		return nil, err
	}
	hashByte := byte(hash(idBytes))
	idBytes = append(idBytes, hashByte)
	return ID(idBytes), nil
}

type ID []byte

func ValidID(id ID) bool {
	return NewService().ValidID(id)
}
func (svc *Service) ValidID(id ID) bool {
	if len(id) != svc.IDLength {
		return false
	}
	hsh := hash([]byte(id)[:svc.IDLength-1])
	return hsh == []byte(id)[svc.IDLength-1]
}

func (id ID) Base64() string {
	return base64.RawURLEncoding.EncodeToString(id)
}

func FromBase64(str string) (ID, error) {
	b, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return ID(b), nil
}
