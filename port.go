package port

import (
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"syscall"
)

var (
	used = map[int]struct{}{}
	m    sync.Mutex
)

func New() int {
	m.Lock()
	defer m.Unlock()

	// iterate through all possible ports, starting from the highest.
	for i := math.MaxUint16; i > 0; i-- {
		// skip if already being used by us.
		if _, ok := used[i]; ok {
			continue
		}

		// connect to the port on any interface.
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", i))
		if err != nil {
			// ensure the syscall error is "connection refused".
			if !errors.Is(err, syscall.ECONNREFUSED) {
				continue
			}

			// set the port as used.
			used[i] = struct{}{}

			return i
		}

		// close the connection.
		conn.Close()
	}

	panic("port exhaustion")
}
