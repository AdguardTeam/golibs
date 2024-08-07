package container_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/container"
)

func ExampleKeyValues() {
	type service struct{}

	var (
		serviceA = service{}
		serviceB = service{}
		serviceC = service{}
	)

	type svcName string

	validate := func(_ service) (err error) {
		return nil
	}

	// These should be validated in this order.
	services := container.KeyValues[svcName, service]{{
		Key:   "svc_a",
		Value: serviceA,
	}, {
		Key:   "svc_b",
		Value: serviceB,
	}, {
		Key:   "svc_c",
		Value: serviceC,
	}}

	for _, kv := range services {
		fmt.Printf("validating service %q: %v\n", kv.Key, validate(kv.Value))
	}

	// Output:
	// validating service "svc_a": <nil>
	// validating service "svc_b": <nil>
	// validating service "svc_c": <nil>
}
