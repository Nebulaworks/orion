package applicant

import "testing"

func TestLockForSameKeyIsSame(t *testing.T) {
	vendor := NewLockVendor()
	firstTime := vendor.LockForName("testing")
	secondTime := vendor.LockForName("testing")
	if firstTime != secondTime {
		t.Fatal()
	}
}

func TestLockForDifferentKeyIsDifferent(t *testing.T) {
	vendor := NewLockVendor()
	firstTime := vendor.LockForName("testing")
	secondTime := vendor.LockForName("different")
	if firstTime == secondTime {
		t.Fatal()
	}
}
