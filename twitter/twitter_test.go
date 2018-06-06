package twitter

import "testing"

func TestNewClient(t *testing.T) {
	tests := []struct {
		a1, a2, a3, a4 string
		wantError      bool
	}{
		{
			"test", "test", "test", "test", false,
		},
		{
			"", "", "", "", true,
		},
	}

	for _, te := range tests {
		_, err := NewClient(te.a1, te.a2, te.a3, te.a4)
		if !te.wantError && err != nil {
			t.Errorf("want no error happen, got an error: %s", err)
		}
		if te.wantError && err == nil {
			t.Error("want error happen, got nothing")
		}
	}
}
