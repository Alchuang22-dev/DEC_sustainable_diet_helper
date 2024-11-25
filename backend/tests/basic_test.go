// tests/basic_test.go
package tests

import (
    "testing"
    "fmt"
)

func TestBasic(t *testing.T) {
    fmt.Println("Running basic test")
    t.Log("This is a test log message")
}