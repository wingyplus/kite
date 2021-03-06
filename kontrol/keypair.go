package kontrol

import (
	"errors"
	"fmt"

	"github.com/koding/cache"
)

// KeyPair defines a single key pair entity
type KeyPair struct {
	// ID is the unique id defining the key pair
	ID string

	// Public key is used to validate tokens
	Public string

	// Private key is used to sign/generate tokens
	Private string
}

func (k *KeyPair) Validate() error {
	if k.ID == "" {
		return errors.New("KeyPair ID field is empty")
	}

	if k.Public == "" {
		return errors.New("KeyPair Public field is empty")
	}

	if k.Private == "" {
		return errors.New("KeyPair Private field is empty")
	}
	return nil
}

// KeyPairStorage is responsible of managing key pairs
type KeyPairStorage interface {
	// AddKey adds the given key pair to the storage
	AddKey(*KeyPair) error

	// DeleteKey deletes the given key pairs from the storage
	DeleteKey(*KeyPair) error

	// GetKeyFromID retrieves the KeyPair from the given ID
	GetKeyFromID(id string) (*KeyPair, error)

	// GetKeyFromPublic retrieves the KeyPairs from the given public Key
	GetKeyFromPublic(publicKey string) (*KeyPair, error)

	// Is valid checks if the given publicKey is valid or not. It's up to the
	// implementer how to implement it. A valid public key returns a nil error.
	IsValid(publicKey string) error
}

func NewMemKeyPairStorage() *MemKeyPairStorage {
	return &MemKeyPairStorage{
		id:     cache.NewMemory(),
		public: cache.NewMemory(),
	}
}

type MemKeyPairStorage struct {
	id     cache.Cache
	public cache.Cache
}

func (m *MemKeyPairStorage) AddKey(keyPair *KeyPair) error {
	if err := keyPair.Validate(); err != nil {
		return err
	}

	m.id.Set(keyPair.ID, keyPair)
	m.public.Set(keyPair.Public, keyPair)
	return nil
}

func (m *MemKeyPairStorage) DeleteKey(keyPair *KeyPair) error {
	if keyPair.Public == "" {
		k, err := m.GetKeyFromID(keyPair.ID)
		if err != nil {
			return err
		}

		m.public.Delete(k.Public)
	}

	m.id.Delete(keyPair.ID)
	return nil
}

func (m *MemKeyPairStorage) GetKeyFromID(id string) (*KeyPair, error) {
	v, err := m.id.Get(id)
	if err != nil {
		return nil, err
	}

	keyPair, ok := v.(*KeyPair)
	if !ok {
		return nil, fmt.Errorf("MemKeyPairStorage: GetKeyFromID value is malformed %+v", v)
	}

	return keyPair, nil
}

func (m *MemKeyPairStorage) GetKeyFromPublic(public string) (*KeyPair, error) {
	v, err := m.public.Get(public)
	if err != nil {
		return nil, err
	}

	keyPair, ok := v.(*KeyPair)
	if !ok {
		return nil, fmt.Errorf("MemKeyPairStorage: GetKeyFromPublic value is malformed %+v", v)
	}

	return keyPair, nil
}

func (m *MemKeyPairStorage) IsValid(public string) error {
	_, err := m.GetKeyFromPublic(public)
	return err
}
