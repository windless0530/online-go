package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNoState signals the persistence file does not exist yet (first launch).
var ErrNoState = errors.New("no persisted state")

// SaveState writes the snapshot atomically: write to a tmp file in the same
// directory, then rename. No fsync — if power dies mid-write the rename is
// either visible or not.
func SaveState(path string, s GameState) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".state-*.tmp")
	if err != nil {
		return fmt.Errorf("create tmp: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close tmp: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

// LoadState reads and validates the snapshot. Returns ErrNoState on first
// launch. On corruption, backs up the bad file to <path>.bak and returns
// ErrNoState so the caller can start fresh.
func LoadState(path string) (GameState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return GameState{}, ErrNoState
		}
		return GameState{}, fmt.Errorf("read state: %w", err)
	}

	var s GameState
	if err := json.Unmarshal(data, &s); err != nil {
		backupCorrupt(path)
		return GameState{}, ErrNoState
	}
	if err := validateState(s); err != nil {
		backupCorrupt(path)
		return GameState{}, ErrNoState
	}
	return s, nil
}

func backupCorrupt(path string) {
	_ = os.Rename(path, path+".bak")
}

func validateState(s GameState) error {
	if len(s.Board) != BoardSize {
		return fmt.Errorf("board size %d != %d", len(s.Board), BoardSize)
	}
	for x, col := range s.Board {
		if len(col) != BoardSize {
			return fmt.Errorf("column %d size %d != %d", x, len(col), BoardSize)
		}
		for y, c := range col {
			if c != Empty && c != Black && c != White {
				return fmt.Errorf("invalid stone %d at (%d,%d)", c, x, y)
			}
		}
	}
	if s.BlackCaptures < 0 || s.WhiteCaptures < 0 {
		return errors.New("negative capture count")
	}
	return nil
}