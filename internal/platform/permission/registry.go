// Package permission keeps the single source of truth for API permission
// names. Protected routes register their name while the route table is built,
// and the seed command reads the same list, so the permissions table can never
// drift from the code.
package permission

import (
	"sort"
	"strings"
	"sync"
)

var registry = struct {
	mu    sync.Mutex
	names []string
}{}

// Register records a permission name. One call happens per protected route at
// registration time.
func Register(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.names = append(registry.names, name)
}

// Reset clears the registry. It runs before a route table is built so repeated
// builds in the same process (console commands, tests) start clean.
func Reset() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.names = nil
}

// All returns every registration in order, duplicates included. A duplicate
// means one permission is mounted on more than one route and must be rejected
// by the route audit test.
func All() []string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	names := make([]string, len(registry.names))
	copy(names, registry.names)
	return names
}

// Names returns the unique permission names, sorted. This is the list the seed
// command writes to the permissions table.
func Names() []string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	seen := make(map[string]struct{}, len(registry.names))
	for _, name := range registry.names {
		seen[name] = struct{}{}
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
