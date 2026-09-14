package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// BackupLedger makes a dated, atomic copy of the ledger at srcPath under dir.
//
// The source must exist and parse (LoadShelf): garbage is never backed up.
// The copy is written atomically (temp file in dir + fsync + rename) with
// permissions 0600. The main ledger is never written by this path.
//
// The destination is <dir>/ledger-YYYYMMDD-HHMMSS.json (UTC timestamp of
// now); if it already exists (two backups in the same second) a -1, -2, …
// suffix is appended. dir is created with MkdirAll if missing.
//
// Returns the path of the backup written.
func BackupLedger(srcPath, dir string, now time.Time) (string, error) {
	// The source must exist (LoadShelf treats a missing ledger as an empty
	// shelf, but backing up nothing is an error) and must parse.
	if _, err := os.Stat(srcPath); err != nil {
		return "", fmt.Errorf("ledger inexistente %s: %w", srcPath, err)
	}
	shelf, err := LoadShelf(srcPath)
	if err != nil {
		return "", err
	}
	data, err := shelf.Marshal()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("creando directorio de backups %s: %w", dir, err)
	}

	dst := uniqueDestination(filepath.Join(dir, now.UTC().Format("ledger-20060102-150405.json")))
	if err := writeAtomic(dst, data); err != nil {
		return "", err
	}
	return dst, nil
}

// PruneBackups removes the oldest ledger-*.json backups under dir, keeping
// the keep newest. keep <= 0 means no limit (no pruning, no-op).
//
// Ordering is by file name, which is chronological because the timestamp
// format (ledger-YYYYMMDD-HHMMSS.json) sorts lexicographically; suffixes
// (-1, -2 from a same-second collision) also sort after their base name.
//
// Invariants: pruning never touches anything outside dir, never removes a
// file that does not match ledger-*.json (the main ledger is never under
// this pattern), and never removes any of the keep newest backups.
// A missing dir is a documented no-op (nothing to prune, no error).
// Returns the removed paths in oldest-first order.
func PruneBackups(dir string, keep int) ([]string, error) {
	if keep <= 0 {
		return nil, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // dir inexistente: nada que podar
		}
		return nil, fmt.Errorf("leyendo directorio de backups %s: %w", dir, err)
	}
	var backups []string
	for _, e := range entries {
		if e.IsDir() || !isLedgerBackupName(e.Name()) {
			continue
		}
		backups = append(backups, e.Name())
	}
	if len(backups) <= keep {
		return nil, nil
	}
	sort.Strings(backups) // nombre == cronológico
	obsolete := backups[:len(backups)-keep]
	var removed []string
	for _, name := range obsolete {
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			return removed, fmt.Errorf("eliminando backup viejo %s: %w", name, err)
		}
		removed = append(removed, filepath.Join(dir, name))
	}
	return removed, nil
}

// isLedgerBackupName reports whether name matches ledger-*.json exactly
// (no subdirectories, nothing else is ever a candidate for pruning).
func isLedgerBackupName(name string) bool {
	return strings.HasPrefix(name, "ledger-") && strings.HasSuffix(name, ".json") &&
		len(name) > len("ledger-.json")
}

// uniqueDestination appends -1, -2, … before the extension until the path
// does not exist (handles two backups within the same second).
func uniqueDestination(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

// writeAtomic writes data to path atomically: temp file in the same
// directory, fsync, chmod 0600, rename over the target.
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".estanteria-backup-*.tmp")
	if err != nil {
		return fmt.Errorf("creando temp: %w", err)
	}
	tmp := f.Name()
	defer os.Remove(tmp) // no-op after a successful rename
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("escribiendo temp: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("fsync temp: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("cerrando temp: %w", err)
	}
	if err := os.Chmod(tmp, 0o600); err != nil {
		return fmt.Errorf("chmod temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("renombrando sobre %s: %w", path, err)
	}
	return nil
}
