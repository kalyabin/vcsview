package vcsview

import "testing"

func TestBranch_Id(t *testing.T) {
	expectedId := "testing_branch"

	b := Branch{}
	b.id = expectedId

	if id := b.Id(); id != expectedId {
		t.Errorf("Branch.Id() = %v, want: %v", id, expectedId)
	}
}

func TestBranch_Head(t *testing.T) {
	expectedHead := "testing_commit"

	b := Branch{}
	b.head = expectedHead

	if head := b.Head(); head != expectedHead {
		t.Errorf("Branch.Head() = %v, want: %v", head, expectedHead)
	}
}

func TestBranch_IsCurrent(t *testing.T) {
	b := Branch{}
	b.isCurrent = true

	if !b.IsCurrent() {
		t.Errorf("Branch.IsCurrent() = false, want: true")
	}
}