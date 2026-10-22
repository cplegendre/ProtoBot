package golden_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRemoveAllRetryRemovesTree(t *testing.T) {
	work := filepath.Join(t.TempDir(), "work")
	gitDir := filepath.Join(work, "clone", ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := removeAllRetry(work); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(work); !os.IsNotExist(err) {
		t.Fatalf("work directory still exists: %v", err)
	}
}

func TestRemoveAllRetryIgnoresMissingPath(t *testing.T) {
	if err := removeAllRetry(filepath.Join(t.TempDir(), "missing")); err != nil {
		t.Fatal(err)
	}
}

func TestReapServerWaitsForExitedProcess(t *testing.T) {
	cmd := exec.Command("true")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	reapServer(cmd)
	if time.Since(start) > time.Second {
		t.Fatal("reapServer blocked on an already-exiting process")
	}
	if cmd.ProcessState == nil {
		t.Fatal("reapServer did not wait for the process")
	}
	reapServer(cmd)
}

func TestReapServerKillsHungProcess(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	reapServer(cmd)
	elapsed := time.Since(start)
	if elapsed > serverWaitTimeout+time.Second {
		t.Fatalf("reapServer took %s", elapsed)
	}
	if cmd.ProcessState == nil {
		t.Fatal("reapServer did not reap the killed process")
	}
}

func TestRemoveAllRetrySucceedsAfterWriterExits(t *testing.T) {
	work := filepath.Join(t.TempDir(), "work")
	gitDir := filepath.Join(work, "clone", ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		lock := filepath.Join(gitDir, "lock")
		for {
			select {
			case <-stop:
				return
			default:
				_ = os.WriteFile(lock, []byte("x"), 0o644)
			}
		}
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
	}()
	if err := removeAllRetry(work); err != nil {
		t.Fatal(err)
	}
	<-done
}

func TestCloseSessionReapsServerBeforeReturn(t *testing.T) {
	cmd := exec.Command("true")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	d := &driver{cmd: cmd, servers: []*exec.Cmd{cmd}}
	d.closeSession()
	if d.cmd != nil {
		t.Fatal("closeSession left cmd set")
	}
	if len(d.servers) != 0 {
		t.Fatal("closeSession left servers set")
	}
	if cmd.ProcessState == nil {
		t.Fatal("closeSession did not wait for the server process")
	}
}

func TestFixtureGitconfigDisablesAutoDetach(t *testing.T) {
	if !strings.Contains(fixtureGitconfig, "autoDetach = false") {
		t.Fatalf("fixture gitconfig does not disable gc.autoDetach:\n%s", fixtureGitconfig)
	}
}
