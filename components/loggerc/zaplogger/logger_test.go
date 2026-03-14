package zaplogger

import "testing"

func TestZapLogger(t *testing.T) {
	_, err := NewZapLogger()

	if err != nil {
		t.Fatalf("something wrong %v", err)
	}
}
