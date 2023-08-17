package edge

import (
	"github.com/robertkrimen/otto"
	"testing"
)

func TestEdge(t *testing.T) {
	vm := otto.New()
	err := vm.Set("input", map[string]any{
		"test": "test",
	})
	v, err := vm.Object("({ test: input.test, test2: input.test })")
	if err != nil {
		t.Fatal(err)
	}
	println(v)
}
