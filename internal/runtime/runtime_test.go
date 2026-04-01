package runtime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTurnControllerActivateCancelRelease(t *testing.T) {
	controller := NewTurnController()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := controller.Activate("th-1", "ses-1", "tu-1", cancel); err != nil {
		t.Fatalf("Activate() unexpected error: %v", err)
	}
	if !controller.IsThreadActive("th-1") {
		t.Fatalf("thread should be active")
	}
	if !controller.HasActiveSession("th-1") {
		t.Fatalf("thread should report active session")
	}
	if !controller.IsSessionActive("th-1", "ses-1") {
		t.Fatalf("session should be active")
	}

	if err := controller.Activate("th-1", "ses-1", "tu-2", cancel); !errors.Is(err, ErrActiveTurnExists) {
		t.Fatalf("second Activate() error = %v, want %v", err, ErrActiveTurnExists)
	}
	if err := controller.Activate("th-1", "ses-2", "tu-2", cancel); err != nil {
		t.Fatalf("Activate(other session) unexpected error: %v", err)
	}

	if err := controller.Cancel("tu-1"); err != nil {
		t.Fatalf("Cancel() unexpected error: %v", err)
	}

	controller.Release("th-1", "ses-1", "tu-1")
	if !controller.IsThreadActive("th-1") {
		t.Fatalf("thread should still be active while another session is running")
	}
	if controller.IsSessionActive("th-1", "ses-1") {
		t.Fatalf("session should be inactive after release")
	}

	if err := controller.Cancel("tu-1"); !errors.Is(err, ErrTurnNotActive) {
		t.Fatalf("Cancel() after release error = %v, want %v", err, ErrTurnNotActive)
	}

	controller.Release("th-1", "ses-2", "tu-2")
	if controller.IsThreadActive("th-1") {
		t.Fatalf("thread should be inactive after releasing all sessions")
	}
	if controller.HasActiveSession("th-1") {
		t.Fatalf("thread should not report active session after releasing all sessions")
	}
}

func TestTurnControllerWaitForIdleAndCancelAll(t *testing.T) {
	controller := NewTurnController()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := controller.Activate("th-1", "ses-1", "tu-1", cancel); err != nil {
		t.Fatalf("Activate() unexpected error: %v", err)
	}
	if got := controller.ActiveCount(); got != 1 {
		t.Fatalf("ActiveCount() = %d, want 1", got)
	}

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer waitCancel()
	if err := controller.WaitForIdle(waitCtx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitForIdle() error = %v, want %v", err, context.DeadlineExceeded)
	}

	cancelled := controller.CancelAll()
	if cancelled != 1 {
		t.Fatalf("CancelAll() = %d, want 1", cancelled)
	}

	controller.Release("th-1", "ses-1", "tu-1")
	waitCtx2, waitCancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer waitCancel2()
	if err := controller.WaitForIdle(waitCtx2); err != nil {
		t.Fatalf("WaitForIdle() after release unexpected error: %v", err)
	}
}

func TestTurnControllerBindTurnSession(t *testing.T) {
	controller := NewTurnController()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := controller.Activate("th-1", "", "tu-1", cancel); err != nil {
		t.Fatalf("Activate() unexpected error: %v", err)
	}
	if err := controller.BindTurnSession("tu-1", "ses-1"); err != nil {
		t.Fatalf("BindTurnSession() unexpected error: %v", err)
	}
	if controller.IsSessionActive("th-1", "") {
		t.Fatalf("empty session scope should be inactive after bind")
	}
	if !controller.IsSessionActive("th-1", "ses-1") {
		t.Fatalf("bound session scope should be active")
	}

	if err := controller.Activate("th-1", "ses-1", "tu-2", cancel); !errors.Is(err, ErrActiveTurnExists) {
		t.Fatalf("Activate(bound session) error = %v, want %v", err, ErrActiveTurnExists)
	}

	controller.Release("th-1", "ses-1", "tu-1")
}

func TestTurnControllerActivateThreadExclusive(t *testing.T) {
	controller := NewTurnController()

	if err := controller.ActivateThreadExclusive("th-1", "guard-1", nil); err != nil {
		t.Fatalf("ActivateThreadExclusive() unexpected error: %v", err)
	}
	if !controller.IsThreadActive("th-1") {
		t.Fatalf("thread should be active while exclusive guard is held")
	}
	if controller.HasActiveSession("th-1") {
		t.Fatalf("exclusive guard should not count as active session")
	}
	if err := controller.Activate("th-1", "ses-1", "tu-1", nil); !errors.Is(err, ErrActiveTurnExists) {
		t.Fatalf("Activate() while thread guard held error = %v, want %v", err, ErrActiveTurnExists)
	}

	controller.ReleaseThreadExclusive("th-1", "guard-1")
	if controller.IsThreadActive("th-1") {
		t.Fatalf("thread should be inactive after releasing exclusive guard")
	}
}
