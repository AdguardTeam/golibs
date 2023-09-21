//go:build windows

package sysresolv

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/log"
	"golang.org/x/exp/slices"
)

// systemResolvers implementation differs for Windows since Go's resolver
// doesn't work there.
//
// See https://github.com/golang/go/issues/33097.
//
// TODO(e.burkov):  See if it's fixed.
type systemResolvers struct {
	// addrs is the slice of cached local resolvers' addresses.
	addrs     []string
	addrsLock sync.RWMutex
}

// newSystemResolvers returns a system resolvers checker for Windows.
func newSystemResolvers(_ HostGenFunc) (r Resolvers) {
	return &systemResolvers{}
}

// type check
var _ Resolvers = (*systemResolvers)(nil)

// Addrs implements the [Resolvers] interface for *systemResolvers.
func (sr *systemResolvers) Addrs() (addrs []string) {
	sr.addrsLock.RLock()
	defer sr.addrsLock.RUnlock()

	return slices.Clone(sr.addrs)
}

// writeExit writes "exit" to w and closes it.  It is supposed to be run in
// a goroutine.
func writeExit(w io.WriteCloser) {
	defer log.OnPanic("systemResolvers: writeExit")

	defer func() {
		closeErr := w.Close()
		if closeErr != nil {
			// TODO(e.burkov):  Remove logging and send the error to a goroutine
			// instead.
			log.Error("systemResolvers: writeExit: closing: %s", closeErr)
		}
	}()

	_, err := io.WriteString(w, "exit")
	if err != nil {
		log.Error("systemResolvers: writeExit: writing: %s", err)
	}
}

// scanAddrs scans the DNS addresses from nslookup's output.  The expected
// output of nslookup looks like this:
//
//	Default Server:  192-168-1-1.qualified.domain.ru
//	Address:  192.168.1.1
func scanAddrs(s *bufio.Scanner) (addrs []string, err error) {
	var errs []error

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[0] != "Address:" {
			continue
		}

		// If the address contains port then it is separated with '#'.
		ipPort := strings.Split(fields[1], "#")
		if len(ipPort) == 0 {
			continue
		}

		addr := ipPort[0]
		if net.ParseIP(addr) == nil {
			errs = append(errs, fmt.Errorf("%q is not a valid ip", addr))

			continue
		}

		addrs = append(addrs, addr)
	}

	if len(addrs) == 0 {
		return nil, errors.List("no addresses", errs...)
	}

	return addrs, nil
}

// getAddrs gets local resolvers' addresses from OS in a special Windows way.
//
// TODO(e.burkov): This whole function needs more detailed research on getting
// local resolvers addresses on Windows.  We execute the external command for
// now that is not the most accurate way.
func (sr *systemResolvers) getAddrs() (addrs []string, err error) {
	cmdPath, err := exec.LookPath("nslookup.exe")
	if err != nil {
		return nil, fmt.Errorf("looking up cmd path: %w", err)
	}

	cmd := exec.Command(cmdPath)

	var stdin io.WriteCloser
	stdin, err = cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("getting the command's stdin pipe: %w", err)
	}

	var stdout io.ReadCloser
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("getting the command's stdout pipe: %w", err)
	}

	go writeExit(stdin)

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start command executing: %w", err)
	}

	s := bufio.NewScanner(stdout)
	addrs, err = scanAddrs(s)
	if err != nil {
		return nil, fmt.Errorf("scan addresses: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("executing the command: %w", err)
	}

	err = s.Err()
	if err != nil {
		return nil, fmt.Errorf("scanning output: %w", err)
	}

	// Don't close StdoutPipe since Wait do it for us in ¿most? cases.
	//
	// See go doc os/exec.Cmd.StdoutPipe.

	return addrs, nil
}

// Refresh implements the [Resolvers] interface for *systemResolvers.
func (sr *systemResolvers) Refresh() (err error) {
	defer func() { err = errors.Annotate(err, "systemResolvers: %w") }()

	got, err := sr.getAddrs()
	if err != nil {
		return fmt.Errorf("can't get addresses: %w", err)
	}
	if len(got) == 0 {
		return nil
	}

	sr.addrsLock.Lock()
	defer sr.addrsLock.Unlock()

	sr.addrs = got

	return nil
}
