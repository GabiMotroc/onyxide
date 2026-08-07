package app

import (
	"bytes"
	"strings"
	"testing"
)

func runAppCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	var out, errBuf bytes.Buffer
	Command.SetArgs(args)
	Command.SetOut(&out)
	Command.SetErr(&errBuf)
	Command.SetIn(strings.NewReader(""))
	err := Command.Execute()
	return out.String(), errBuf.String(), err
}

func TestCLIAdd(t *testing.T) {
	isolatedDataDir(t)

	if _, _, err := runAppCommand(t, "add", "rider"); err != nil {
		t.Fatalf("app add rider: %v", err)
	}

	apps, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(apps) != 1 || apps[0].Name != "rider" {
		t.Errorf("LoadApps() = %v, want [{rider false}]", apps)
	}
}

func TestCLIAddDuplicate(t *testing.T) {
	isolatedDataDir(t)

	if _, _, err := runAppCommand(t, "add", "rider"); err != nil {
		t.Fatalf("first app add rider: %v", err)
	}

	_, _, err := runAppCommand(t, "add", "RIDER")
	if err == nil {
		t.Fatal("app add RIDER: error = nil, want duplicate error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("app add RIDER error = %v, want contains %q", err, "already exists")
	}
}

func TestCLIList(t *testing.T) {
	isolatedDataDir(t)

	if _, _, err := runAppCommand(t, "add", "rider"); err != nil {
		t.Fatalf("app add rider: %v", err)
	}
	if _, _, err := runAppCommand(t, "add", "code"); err != nil {
		t.Fatalf("app add code: %v", err)
	}

	out, _, err := runAppCommand(t, "list")
	if err != nil {
		t.Fatalf("app list: %v", err)
	}
	if !strings.Contains(out, "rider") || !strings.Contains(out, "code") {
		t.Errorf("app list output = %q, want rider and code", out)
	}
}

func TestCLIListEmpty(t *testing.T) {
	isolatedDataDir(t)

	out, _, err := runAppCommand(t, "list")
	if err != nil {
		t.Fatalf("app list: %v", err)
	}
	if strings.Contains(out, "rider") {
		t.Errorf("app list output = %q, want no entries", out)
	}
}

func TestCLIRemove(t *testing.T) {
	isolatedDataDir(t)

	for _, name := range []string{"rider", "code"} {
		if _, _, err := runAppCommand(t, "add", name); err != nil {
			t.Fatalf("app add %s: %v", name, err)
		}
	}

	out, _, err := runAppCommand(t, "remove", "RIDER")
	if err != nil {
		t.Fatalf("app remove RIDER: %v", err)
	}
	if !strings.Contains(out, `app "RIDER" removed`) {
		t.Errorf("app remove output = %q, want %q", out, `app "RIDER" removed`)
	}

	apps, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(apps) != 1 || apps[0].Name != "code" {
		t.Errorf("LoadApps() after remove = %v, want [{code false}]", apps)
	}
}

func TestCLIRemoveNotFound(t *testing.T) {
	isolatedDataDir(t)

	_, _, err := runAppCommand(t, "remove", "vim")
	if err == nil {
		t.Fatal("app remove vim: error = nil, want not-found error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("app remove vim error = %v, want contains %q", err, "not found")
	}
}

func TestCLIClear(t *testing.T) {
	isolatedDataDir(t)

	for _, name := range []string{"rider", "code"} {
		if _, _, err := runAppCommand(t, "add", name); err != nil {
			t.Fatalf("app add %s: %v", name, err)
		}
	}

	if _, _, err := runAppCommand(t, "clear"); err != nil {
		t.Fatalf("app clear: %v", err)
	}

	apps, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(apps) != 0 {
		t.Errorf("LoadApps() after clear = %v, want empty", apps)
	}
}

func TestCLIWorkflow(t *testing.T) {
	isolatedDataDir(t)

	steps := []struct {
		args  []string
		wantErr bool
	}{
		{[]string{"add", "rider"}, false},
		{[]string{"add", "rider"}, true},
		{[]string{"add", "code"}, false},
		{[]string{"list"}, false},
		{[]string{"remove", "code"}, false},
		{[]string{"remove", "code"}, true},
		{[]string{"clear"}, false},
	}
	for _, step := range steps {
		_, _, err := runAppCommand(t, step.args...)
		if (err != nil) != step.wantErr {
			t.Errorf("app %v: err = %v, wantErr = %v", step.args, err, step.wantErr)
		}
	}

	apps, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(apps) != 0 {
		t.Errorf("LoadApps() after workflow = %v, want empty", apps)
	}
}
