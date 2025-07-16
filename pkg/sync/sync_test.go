package main

import (
	"testing"
	"time"
)

func TestNewCloser(t *testing.T) {
	closer := NewCloser()
	
	if closer == nil {
		t.Fatal("NewCloser() returned nil")
	}
	
	if closer.doneCh == nil {
		t.Fatal("NewCloser() did not initialize doneCh")
	}
}

func TestCloser_Done(t *testing.T) {
	closer := NewCloser()
	
	done := closer.Done()
	if done == nil {
		t.Fatal("Done() returned nil channel")
	}
	
	// Channel should not be closed initially
	select {
	case <-done:
		t.Fatal("Done() channel should not be closed initially")
	default:
		// Expected behavior - channel is open
	}
}

func TestCloser_Close(t *testing.T) {
	closer := NewCloser()
	
	// Close the closer
	closer.Close()
	
	// Channel should be closed after Close()
	select {
	case <-closer.Done():
		// Expected behavior - channel is closed
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Done() channel should be closed after Close()")
	}
}

func TestCloser_CloseMultipleTimes(t *testing.T) {
	closer := NewCloser()
	
	// Close multiple times - should not panic
	closer.Close()
	closer.Close()
	closer.Close()
	
	// Channel should still be properly closed
	select {
	case <-closer.Done():
		// Expected behavior - channel is closed
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Done() channel should be closed after multiple Close() calls")
	}
}

func TestCloser_CloseRace(t *testing.T) {
	closer := NewCloser()
	
	// Test concurrent Close() calls
	done := make(chan struct{}, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			closer.Close()
			done <- struct{}{}
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Channel should be closed
	select {
	case <-closer.Done():
		// Expected behavior - channel is closed
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Done() channel should be closed after concurrent Close() calls")
	}
}
