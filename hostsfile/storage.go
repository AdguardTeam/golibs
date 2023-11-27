package hostsfile

import (
	"fmt"
	"io"
	"net/netip"
	"strings"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/log"
	"golang.org/x/exp/slices"
)

// Storage indexes the hosts file records.
type Storage interface {
	// ByAddr returns the hostnames for the given address.
	ByAddr(addr netip.Addr) (names []string)

	// ByName returns the addresses for the given hostname.
	ByName(name string) (addrs []netip.Addr)
}

// unit is a convenient alias for empty struct.
type unit = struct{}

// set is a helper type that removes duplicates.
//
// TODO(e.burkov):  Think of using slices in combination with binary search.
type set[K string | netip.Addr] map[K]unit

// orderedSet is a helper type for storing values in original adding order and
// dealing with duplicates.
type orderedSet[K string | netip.Addr] struct {
	set  set[K]
	vals []K
}

// add adds val to os if it's not already there.
func (os *orderedSet[K]) add(key, val K) {
	if _, ok := os.set[key]; !ok {
		os.set[key] = unit{}
		os.vals = append(os.vals, val)
	}
}

// Convenience aliases for [orderedSet].
type (
	namesSet = orderedSet[string]
	addrsSet = orderedSet[netip.Addr]
)

// DefaultStorage is a [Storage] that removes duplicates.  It also implements
// the [HandleSet] interface and therefore can be used within [Parse].
//
// It must be initialized with [NewDefaultStorage].
type DefaultStorage struct {
	// names maps each address to its names in original case and in original
	// adding order without duplicates.
	names map[netip.Addr]*namesSet

	// addrs maps each host to its addresses in original adding order without
	// duplicates.
	addrs map[string]*addrsSet
}

// NewDefaultStorage parses data if hosts files format from readers and returns
// a new properly initialized DefaultStorage.  readers are optional, an empty
// storage is completely usable.
func NewDefaultStorage(readers ...io.Reader) (s *DefaultStorage, err error) {
	s = &DefaultStorage{
		names: map[netip.Addr]*namesSet{},
		addrs: map[string]*addrsSet{},
	}

	for i, r := range readers {
		if err = Parse(s, r, nil); err != nil {
			return nil, fmt.Errorf("reader at index %d: %w", i, err)
		}
	}

	return s, nil
}

// type check
var _ HandleSet = (*DefaultStorage)(nil)

// Add implements the [Set] interface for *DefaultStorage.  It skips records
// without hostnames, ignores duplicates and squashes the rest.
func (s *DefaultStorage) Add(rec *Record) {
	names := s.names[rec.Addr]
	if names == nil {
		names = &namesSet{set: set[string]{}}
		s.names[rec.Addr] = names
	}

	for _, name := range rec.Names {
		lowered := strings.ToLower(name)
		names.add(lowered, name)

		addrs := s.addrs[lowered]
		if addrs == nil {
			addrs = &addrsSet{
				vals: []netip.Addr{},
				set:  set[netip.Addr]{},
			}
			s.addrs[lowered] = addrs
		}
		addrs.add(rec.Addr, rec.Addr)
	}
}

// HandleInvalid implements the [HandleSet] interface for *DefaultStorage.  It
// essentially ignores empty lines and logs all other errors at debug level.
func (s *DefaultStorage) HandleInvalid(srcName string, _ []byte, err error) {
	lineErr := &LineError{}
	if !errors.As(err, &lineErr) {
		log.Debug("hostsfile: unexpected parsing error: %s", err)

		return
	}

	if errors.Is(err, ErrEmptyLine) {
		// Ignore empty lines and comments.
		return
	}

	log.Debug("hostsfile: source %q: %s", srcName, lineErr)
}

// type check
var _ Storage = (*DefaultStorage)(nil)

// ByAddr implements the [Storage] interface for *DefaultStorage.  It returns
// each host for addr in original case, in original adding order without
// duplicates.  It returns nil if h doesn't contain the addr.
func (s *DefaultStorage) ByAddr(addr netip.Addr) (hosts []string) {
	if hostsSet, ok := s.names[addr]; ok {
		hosts = hostsSet.vals
	}

	return hosts
}

// ByName implements the [Storage] interface for *DefaultStorage.  It returns
// each address for host in original adding order without duplicates.  It
// returns nil if h doesn't contain the host.
func (s *DefaultStorage) ByName(host string) (addrs []netip.Addr) {
	if addrsSet, ok := s.addrs[strings.ToLower(host)]; ok {
		addrs = addrsSet.vals
	}

	return addrs
}

// RangeNames ranges through all addresses in s and calls f with all the
// corresponding names for each one.  The order of range is undefined.  names
// must not be modified.
func (s *DefaultStorage) RangeNames(f func(addr netip.Addr, names []string) (cont bool)) {
	for addr, names := range s.names {
		if !f(addr, names.vals) {
			return
		}
	}
}

// RangeAddrs ranges through all hostnames in s and calls f with all the
// corresponding addresses for each one.  The order of range is undefined.
// addrs must not be modified.
func (s *DefaultStorage) RangeAddrs(f func(host string, addrs []netip.Addr) (cont bool)) {
	for host, addrs := range s.addrs {
		if !f(host, addrs.vals) {
			return
		}
	}
}

// Equal returns true if s and other contain the same addresses mapped to the
// same hostnames.
func (s *DefaultStorage) Equal(other *DefaultStorage) (ok bool) {
	if s == nil || other == nil {
		return s == other
	} else if len(s.names) != len(other.names) || len(s.addrs) != len(other.addrs) {
		return false
	}

	// Don't use [maps.Equal] here, since we're only interested in comparing the
	// fields of values, not the values itself.
	for addr, names := range s.names {
		var otherNames *namesSet
		switch otherNames, ok = other.names[addr]; {
		case
			!ok,
			len(names.vals) != len(otherNames.vals),
			!slices.Equal(names.vals, otherNames.vals):
			return false
		}
	}

	return true
}
