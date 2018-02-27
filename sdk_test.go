package geetest

import "testing"

func TestNewSDK(t *testing.T) {
    tester := New("ID", "KEY")

    if tester.Register == nil {
        t.Fatal()
    }
}
