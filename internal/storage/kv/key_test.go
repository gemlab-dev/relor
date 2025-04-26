package kv_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/gemlab-dev/relor/internal/storage/kv"
	"github.com/google/uuid"
)

func TestWorkflowKeyOrder(t *testing.T) {
	// Test the WorkflowKey function
	first, err := kv.WorkflowKey([]byte("prefix"), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err := kv.WorkflowKey([]byte("prefix"), uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bytes.Compare(first, second) >= 0 {
		t.Errorf("Expected first key to be less than second key, got %s >= %s", first, second)
	}
}

func TestWorkflowKeyEqual(t *testing.T) {
	// Test the WorkflowKey function
	first, err := kv.WorkflowKey([]byte("prefix"), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err := kv.WorkflowKey([]byte("prefix"), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Errorf("Expected keys to be equal, got %s != %s", first, second)
	}
}

func TestWorkflowKeyNoPrefix(t *testing.T) {
	// Test the WorkflowKey function
	first, err := kv.WorkflowKey(nil, uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err := kv.WorkflowKey(nil, uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bytes.Compare(first, second) >= 0 {
		t.Errorf("Expected first key to be less than second key, got %s >= %s", first, second)
	}
}

func TestTransitionKeyDescendingOrder(t *testing.T) {
	first := kv.TransitionKey([]byte("prefix"), 0)
	second := kv.TransitionKey([]byte("prefix"), 1)
	third := kv.TransitionKey([]byte("prefix"), 1<<64-1)

	if bytes.Compare(first, second) <= 0 {
		t.Errorf("Expected first key to be greater than second key, got %s >= %s", first, second)
	}
	if bytes.Compare(second, third) <= 0 {
		t.Errorf("Expected second key to be greater than third key, got %s >= %s", second, third)
	}
	if bytes.Compare(first, third) <= 0 {
		t.Errorf("Expected first key to be greater than third key, got %s >= %s", first, third)
	}
}

func TestTransitionKeyEmptyPrefix(t *testing.T) {
	first := kv.TransitionKey(nil, 0)
	second := kv.TransitionKey(nil, 1)

	if bytes.Compare(first, second) <= 0 {
		t.Errorf("Expected first key to be greater than second key, got %s >= %s", first, second)
	}
}

func TestScheduleKeyAscendingOrder(t *testing.T) {
	first, err := kv.ScheduleKey([]byte("prefix"), time.Unix(0, 0), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err := kv.ScheduleKey([]byte("prefix"), time.Unix(1, 0), uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bytes.Compare(first, second) >= 0 {
		t.Errorf("Expected first key to be less than second key, got %s >= %s", first, second)
	}

	// Make sure the key IDs don't affect the order.
	first, err = kv.ScheduleKey([]byte("prefix"), time.Unix(0, 0), uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err = kv.ScheduleKey([]byte("prefix"), time.Unix(0, 1), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bytes.Compare(first, second) >= 0 {
		t.Errorf("Expected first key to be less than second key, got %s >= %s", first, second)
	}
}

func TestScheduleKeySameTimeDifferentKeys(t *testing.T) {
	first, err := kv.ScheduleKey([]byte("prefix"), time.Unix(0, 0), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err := kv.ScheduleKey([]byte("prefix"), time.Unix(0, 0), uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bytes.Equal(first, second) {
		t.Errorf("Expected keys to be different, got %s == %s", first, second)
	}
}

func TestScheduleKeyEmptyPrefix(t *testing.T) {
	first, err := kv.ScheduleKey(nil, time.Unix(0, 0), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	second, err := kv.ScheduleKey(nil, time.Unix(1, 0), uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bytes.Compare(first, second) >= 0 {
		t.Errorf("Expected first key to be less than second key, got %s >= %s", first, second)
	}
}
func TestScheduleKeyZeroTime(t *testing.T) {
	_, err := kv.ScheduleKey([]byte("prefix"), time.Time{}, uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
