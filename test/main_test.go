package test

import (
	test_util "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"golang.org/x/net/context"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// running test container
	test_util.SetUpTestContainers(ctx)

	test_util.StartServer()
	code := m.Run()

	// clean
	test_util.Terminate(ctx)

	os.Exit(code)
}
