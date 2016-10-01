package fixedhashmap

import "testing"

func TestHappyPath(t *testing.T) {
	els := map[string]int{
		"sator": 1,
		"arepo": 2,
		"tenet": 3,
		"opera": 4,
		"rotas": 5,
	}
	m := New(5)
	for key, value := range els {
		m.Set(key, value)
	}

	for key, value := range els {
		got := m.Get(key).(int)
		if got != value {
			t.Errorf("Got %v, want %v", got, value)
		}
	}

	for key, _ := range els {
		m.Delete(key)
	}
	load := m.Load()
	if load != 0 {
		t.Errorf("Unexpected load: got %v, want %v", load, 0)
	}
}

// TestSetMaxLoad tests that we do not overflow if setting too many elements.
func TestSetMaxLoad(t *testing.T) {
	m := New(2)
	m.Set("a", 1)
	m.Set("b", 2)
	if m.Set("c", 3) {
		t.Errorf("Expected set of c to be false")
	}
	if m.Get("a") != 1 {
		t.Errorf("Expected a to be 1")
	}
	if m.Get("b") != 2 {
		t.Errorf("Expected b to be 2")
	}
	if m.Get("c") != nil {
		t.Errorf("Expected c to be nil")
	}
	if m.Load() != 1 {
		t.Errorf("Expected load to be 1")
	}
}

func TestDelete(t *testing.T) {
	m := New(2)
	m.Set("a", 1)
	m.Set("b", 2)
	m.Delete("b")

	if m.Get("b") != nil {
		t.Errorf("Expected get b to be nil")
	}
	m.Set("b", 5)
	got := m.Get("b")
	if got != 5 {
		t.Errorf("Expected 5, got %v", got)
	}
}
