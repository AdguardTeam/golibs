package hostsfile_test

import (
	"bytes"
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/hostsfile"
	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleFuncSet() {
	const content = "# comment\n" +
		"1.2.3.4 host1 host2\n" +
		"4.3.2.1 host3\n" +
		"1.2.3.4 host4 host5 # repeating address\n" +
		"2.3.4.5 host3 # repeating hostname"

	addrs := map[string][]netip.Addr{}
	names := map[netip.Addr][]string{}
	set := hostsfile.FuncSet(func(_ context.Context, r *hostsfile.Record) {
		names[r.Addr] = append(names[r.Addr], r.Names...)
		for _, name := range r.Names {
			addrs[name] = append(addrs[name], r.Addr)
		}
	})

	// Parse the hosts file.
	ctx := context.Background()
	err := hostsfile.Parse(ctx, set, strings.NewReader(content), nil)
	fmt.Printf("error: %s\n", err)
	fmt.Printf("records for 1.2.3.4: %q\n", names[netip.MustParseAddr("1.2.3.4")])
	fmt.Printf("records for host3: %s\n", addrs["host3"])

	// Output:
	// error: parsing: line 1: line is empty
	// records for 1.2.3.4: ["host1" "host2" "host4" "host5"]
	// records for host3: [4.3.2.1 2.3.4.5]
}

// invalidSet is a [HandleSet] implementation that collects invalid records.
type invalidSet []hostsfile.Record

// type check
var _ hostsfile.HandleSet = (*invalidSet)(nil)

// Add implements the [Set] interface for invalidSet.
func (s *invalidSet) Add(_ context.Context, r *hostsfile.Record) { *s = append(*s, *r) }

// HandleInvalid implements the [HandleSet] interface for invalidSet.
func (s *invalidSet) HandleInvalid(ctx context.Context, srcName string, data []byte, err error) {
	addrErr := &netutil.AddrError{}
	if !errors.As(err, &addrErr) {
		return
	}

	rec := &hostsfile.Record{Source: srcName}

	_ = rec.UnmarshalText(data)
	if commIdx := bytes.IndexByte(data, '#'); commIdx >= 0 {
		data = bytes.TrimRight(data[:commIdx], " \t")
	}

	invIdx := bytes.Index(data, []byte(addrErr.Addr))
	for _, name := range bytes.Fields(data[invIdx:]) {
		rec.Names = append(rec.Names, string(name))
	}

	s.Add(ctx, rec)
}

func ExampleHandleSet() {
	const content = "\n" +
		"# comment\n" +
		"4.3.2.1 invalid.-host valid.host # comment\n" +
		"1.2.3.4 another.valid.host\n"

	set := invalidSet{}
	ctx := context.Background()
	err := hostsfile.Parse(ctx, &set, strings.NewReader(content), nil)
	fmt.Printf("error: %v\n", err)
	for _, r := range set {
		fmt.Printf("%q\n", r.Names)
	}

	// Output:
	// error: <nil>
	// ["invalid.-host" "valid.host"]
	// ["another.valid.host"]
}
