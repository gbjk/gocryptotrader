package balance

import (
	"testing"
	"time"
)

func TestBalance_Add(t *testing.T) {
	tests := []struct {
		name      string
		initial   Balance
		add       Balance
		expected  Balance
		expectErr bool
	}{
		{
			name: "simple add",
			initial: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 100,
			},
			add: Balance{
				Free:      5,
				Reserved:  5,
				UpdatedAt: 200,
			},
			expected: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 200,
			},
			expectErr: false,
		},
		{
			name: "add with older timestamp",
			initial: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 200,
			},
			add: Balance{
				Free:      5,
				Reserved:  5,
				UpdatedAt: 100,
			},
			expected: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 200,
			},
			expectErr: false,
		},
		{
			name: "add to zero balance",
			initial: Balance{
				Free:      0,
				Reserved:  0,
				UpdatedAt: 0,
			},
			add: Balance{
				Free:      5,
				Reserved:  5,
				UpdatedAt: 100,
			},
			expected: Balance{
				Free:      5,
				Reserved:  5,
				UpdatedAt: 100,
			},
			expectErr: false,
		},
		{
			name: "add zero balance",
			initial: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 100,
			},
			add: Balance{
				Free:      0,
				Reserved:  0,
				UpdatedAt: 200,
			},
			expected: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 200,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.initial.Add(tt.add)
			if (err != nil) != tt.expectErr {
				t.Errorf("Balance.Add() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if result.Free != tt.expected.Free || result.Reserved != tt.expected.Reserved || result.UpdatedAt != tt.expected.UpdatedAt {
				t.Errorf("Balance.Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProtectedBalance_reset(t *testing.T) {
	pb := &ProtectedBalance{
		balance: Balance{
			Free:      100,
			Reserved:  50,
			UpdatedAt: 1,
		},
	}

	pb.reset()

	if pb.balance.Free != 0 {
		t.Errorf("After reset(), Free = %d, want %d", pb.balance.Free, 0)
	}
	if pb.balance.Reserved != 0 {
		t.Errorf("After reset(), Reserved = %d, want %d", pb.balance.Reserved, 0)
	}
	if pb.balance.UpdatedAt != 0 {
		t.Errorf("After reset(), UpdatedAt = %d, want %d", pb.balance.UpdatedAt, 0)
	}

	// Test resetting an already zero balance
	pb.reset() // Call reset again
	if pb.balance.Free != 0 || pb.balance.Reserved != 0 || pb.balance.UpdatedAt != 0 {
		t.Errorf("After resetting a zero balance, expected all fields to be zero, got %+v", pb.balance)
	}
}

func TestBalance_Wait(t *testing.T) {
	t.Run("waits for change", func(t *testing.T) {
		b := &Balance{Free: 10, UpdatedAt: 100}
		go func() {
			time.Sleep(50 * time.Millisecond)
			b.update(Balance{Free: 20, UpdatedAt: 200})
		}()

		changed, err := b.Wait(100, 150*time.Millisecond)
		if err != nil {
			t.Fatalf("Wait returned an error: %v", err)
		}
		if !changed {
			t.Errorf("Wait returned false, expected true")
		}
		if b.Free != 20 || b.UpdatedAt != 200 {
			t.Errorf("Balance not updated correctly, got %+v", b)
		}
	})

	t.Run("respects maxWait", func(t *testing.T) {
		b := &Balance{Free: 10, UpdatedAt: 100}
		// No update will occur

		changed, err := b.Wait(100, 50*time.Millisecond)
		if err != nil {
			t.Fatalf("Wait returned an error: %v", err)
		}
		if changed {
			t.Errorf("Wait returned true, expected false (timeout)")
		}
		// Balance should remain unchanged
		if b.Free != 10 || b.UpdatedAt != 100 {
			t.Errorf("Balance changed unexpectedly, got %+v", b)
		}
	})

	t.Run("returns immediately if already changed", func(t *testing.T) {
		b := &Balance{Free: 10, UpdatedAt: 200} // UpdatedAt is already newer

		changed, err := b.Wait(100, 50*time.Millisecond)
		if err != nil {
			t.Fatalf("Wait returned an error: %v", err)
		}
		if !changed {
			t.Errorf("Wait returned false, expected true (already changed)")
		}
	})

	t.Run("wait with zero maxWait", func(t *testing.T) {
		b := &Balance{Free: 10, UpdatedAt: 100}

		// Expect to return immediately, no change
		changed, err := b.Wait(100, 0)
		if err != nil {
			t.Fatalf("Wait returned an error: %v", err)
		}
		if changed {
			t.Errorf("Wait returned true, expected false (zero maxWait)")
		}
		if b.Free != 10 || b.UpdatedAt != 100 {
			t.Errorf("Balance changed unexpectedly, got %+v", b)
		}

		// Now test with an already newer timestamp
		b2 := &Balance{Free: 10, UpdatedAt: 200}
		changed2, err2 := b2.Wait(100, 0)
		if err2 != nil {
			t.Fatalf("Wait returned an error: %v", err2)
		}
		if !changed2 {
			t.Errorf("Wait returned false, expected true (zero maxWait, already changed)")
		}
	})
}

func TestProtectedBalance_GetFree(t *testing.T) {
	pb := &ProtectedBalance{
		balance: Balance{
			Free:      100,
			Reserved:  50,
			UpdatedAt: 1,
		},
	}

	freeBalance := pb.GetFree()
	if freeBalance != 100 {
		t.Errorf("GetFree() = %d, want %d", freeBalance, 100)
	}

	// Test with zero balance
	pb.balance.Free = 0
	freeBalance = pb.GetFree()
	if freeBalance != 0 {
		t.Errorf("GetFree() with zero balance = %d, want %d", freeBalance, 0)
	}
}

func TestBalance_update(t *testing.T) {
	tests := []struct {
		name      string
		initial   Balance
		update    Balance
		expected  Balance
		expectErr bool
	}{
		{
			name: "successful update",
			initial: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 100,
			},
			update: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 200,
			},
			expected: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 200,
			},
			expectErr: false,
		},
		{
			name: "out-of-sequence update",
			initial: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 200,
			},
			update: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 100,
			},
			expected: Balance{ // Should not update
				Free:      10,
				Reserved:  5,
				UpdatedAt: 200,
			},
			expectErr: true, // Expect an error for out-of-sequence update
		},
		{
			name: "update with zero UpdatedAt",
			initial: Balance{
				Free:      10,
				Reserved:  5,
				UpdatedAt: 100,
			},
			update: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 0, // Zero UpdatedAt should not be applied
			},
			expected: Balance{ // Should not update
				Free:      10,
				Reserved:  5,
				UpdatedAt: 100,
			},
			expectErr: true, // Expect an error for zero UpdatedAt
		},
		{
			name: "initial update (UpdatedAt is 0)",
			initial: Balance{
				Free:      0,
				Reserved:  0,
				UpdatedAt: 0,
			},
			update: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 100,
			},
			expected: Balance{
				Free:      15,
				Reserved:  10,
				UpdatedAt: 100,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.initial // Use a copy to avoid modifying the initial struct
			err := b.update(tt.update)
			if (err != nil) != tt.expectErr {
				t.Errorf("balance.update() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if b.Free != tt.expected.Free || b.Reserved != tt.expected.Reserved || b.UpdatedAt != tt.expected.UpdatedAt {
				t.Errorf("balance.update() = %v, want %v", b, tt.expected)
			}
		})
	}
}
