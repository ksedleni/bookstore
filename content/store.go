package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Store is a local filesystem content store keyed by content hash.
// Any third party can implement the same layout or serve these JSON envelopes.
type Store struct {
	Root string
}

// NewStore ensures the root directory exists.
func NewStore(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &Store{Root: root}, nil
}

func (s *Store) path(contentHash string) string {
	return filepath.Join(s.Root, contentHash+".json")
}

// Put writes an envelope (overwrites if present).
func (s *Store) Put(env *Envelope) error {
	if env.ContentHash == "" {
		return fmt.Errorf("missing content_hash")
	}
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(env.ContentHash), data, 0o644)
}

// Get loads an envelope by content hash.
func (s *Store) Get(contentHash string) (*Envelope, error) {
	data, err := os.ReadFile(s.path(contentHash))
	if err != nil {
		return nil, err
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return &env, nil
}

// Exists reports whether an envelope is stored.
func (s *Store) Exists(contentHash string) bool {
	_, err := os.Stat(s.path(contentHash))
	return err == nil
}
