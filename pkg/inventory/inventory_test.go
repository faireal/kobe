package inventory

import (
	"fmt"
	"github.com/faireal/kobe/pkg/constant"
	"os"
	"testing"
)

var (
	host = "10.10.10.88"
	port = 8080
)

func TestKobeInventoryProvider_ListHandler(t *testing.T) {
	id := "79aa7416-133f-4678-8425-7d87e6e613ed"
	provider := NewKobeInventoryProvider(host, port)
	os.Setenv(constant.TaskEnvKey, id)
	result, err := provider.ListHandler()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
