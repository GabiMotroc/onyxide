package app

import (
	"os"
	"path/filepath"
	"testing"
)

func isolatedDataDir(t *testing.T) {
	t.Helper()
	t.Setenv("APPDATA", "")
	t.Setenv("HOME", t.TempDir())
}

func TestContainsAppName(t *testing.T) {
	apps := []App{
		{Name: "rider", IsTerminal: false},
		{Name: "Code", IsTerminal: true},
	}

	tests := []struct {
		name string
		apps []App
		needle string
		want  bool
	}{
		{"exact match", apps, "rider", true},
		{"case insensitive", apps, "RIDER", true},
		{"case insensitive reversed", apps, "code", true},
		{"no match", apps, "vim", false},
		{"empty list", []App{}, "rider", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAppName(tt.apps, tt.needle); got != tt.want {
				t.Errorf("ContainsAppName(%v, %q) = %v, want %v", tt.apps, tt.needle, got, tt.want)
			}
		})
	}
}

func TestRemoveByName(t *testing.T) {
	apps := []App{
		{Name: "rider"},
		{Name: "Code"},
		{Name: "rider"},
	}

	filtered, found := RemoveByName(apps, "RIDER")
	if !found {
		t.Fatal("RemoveByName(...) found = false, want true")
	}
	if len(filtered) != 1 {
		t.Fatalf("RemoveByName(...) returned %d apps, want 1: %v", len(filtered), filtered)
	}
	if filtered[0].Name != "Code" {
		t.Errorf("RemoveByName(...) kept %q, want %q", filtered[0].Name, "Code")
	}

	filtered, found = RemoveByName(apps, "vim")
	if found {
		t.Error("RemoveByName(..., \"vim\") found = true, want false")
	}
	if len(filtered) != len(apps) {
		t.Errorf("RemoveByName(..., \"vim\") changed list length: got %d, want %d", len(filtered), len(apps))
	}
}

func TestSaveLoadAppsRoundTrip(t *testing.T) {
	isolatedDataDir(t)
	want := []App{
		{Name: "rider"},
		{Name: "Code", IsTerminal: true},
	}
	if err := SaveApps(want); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	got, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("LoadApps() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("LoadApps()[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestLoadAppsCreatesMissingFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("APPDATA", "")
	t.Setenv("HOME", home)

	apps, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(apps) != 0 {
		t.Fatalf("LoadApps() on missing file = %v, want empty", apps)
	}

	wantPath := filepath.Join(home, ".local", "share", "onyxide", "persistent", "apps.json")
	if _, err := os.Stat(wantPath); err != nil {
		t.Errorf("expected apps.json to be created at %s: %v", wantPath, err)
	}
}

func TestIsTerminal(t *testing.T) {
	isolatedDataDir(t)

	if err := SaveApps([]App{{Name: "code", IsTerminal: true}}); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	tests := []struct {
		name string
		want bool
	}{
		{"code", true},
		{"CODE", true},
		{"rider", false},
	}
	for _, tt := range tests {
		if got := IsTerminal(tt.name); got != tt.want {
			t.Errorf("IsTerminal(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}
