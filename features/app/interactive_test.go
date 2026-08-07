package app

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

// tuiStage is one chunk of a staged keyboard script. wait is the delay to
// sleep after the script bytes are written (used to let a lone ESC resolve
// after ultraviolet's 50ms escape timeout, and to let events settle).
type tuiStage struct {
	script string
	wait   time.Duration
}

// runTUIStaged is like runTUI but feeds the script in stages through an
// io.Pipe, which stays open between writes. This lets tests pace input that a
// single byte stream cannot express unambiguously, e.g. a lone ESC key.
func runTUIStaged(t *testing.T, stages []tuiStage, settle time.Duration) (model, error) {
	t.Helper()
	var out bytes.Buffer
	pr, pw := io.Pipe()
	p := tea.NewProgram(initialModel(), tea.WithInput(pr), tea.WithOutput(&out), tea.WithoutSignals())

	done := make(chan struct{})
	var m tea.Model
	var runErr error
	go func() {
		m, runErr = p.Run()
		close(done)
	}()
	for _, s := range stages {
		if _, err := pw.Write([]byte(s.script)); err != nil {
			t.Fatalf("writing stage %q: %v", s.script, err)
		}
		if s.wait > 0 {
			time.Sleep(s.wait)
		}
	}
	if settle > 0 {
		time.Sleep(settle)
		p.Quit()
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("program did not exit")
	}
	pw.Close()
	if m == nil {
		t.Fatalf("program failed to return a model: %v", runErr)
	}
	return m.(model), runErr
}

// runTUI drives the real Bubbletea program end to end: a scripted byte stream
// is fed through the program's input reader, exercising the actual keyboard
// input pipeline (bytes -> terminal reader -> KeyPressMsg -> model.Update).
//
// If settle > 0 the program is terminated externally via p.Quit() after that
// delay (needed to observe mid-flow state, since some modes have no quit key);
// otherwise the script itself must end with a quit key (q/s).
func runTUI(t *testing.T, script string, settle time.Duration) (model, error) {
	t.Helper()
	var out bytes.Buffer
	var in bytes.Buffer
	in.WriteString(script)
	p := tea.NewProgram(initialModel(), tea.WithInput(&in), tea.WithOutput(&out), tea.WithoutSignals())

	done := make(chan struct{})
	var m tea.Model
	var runErr error
	go func() {
		m, runErr = p.Run()
		close(done)
	}()
	if settle > 0 {
		// Processing the scripted keys is microseconds of work; the margin is
		// orders of magnitude larger. QuitMsg is FIFO behind the keypresses.
		time.Sleep(settle)
		p.Quit()
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("program did not exit; script must end with a quit key or settle > 0")
	}
	if m == nil {
		t.Fatalf("program failed to return a model: %v", runErr)
	}
	return m.(model), runErr
}

func TestTUICommandWiring(t *testing.T) {
	isolatedDataDir(t)

	var out bytes.Buffer
	Command.SetArgs(nil)
	Command.SetIn(strings.NewReader("q"))
	Command.SetOut(&out)
	Command.SetErr(&out)
	if err := Command.Execute(); err != nil {
		t.Fatalf("app command TUI: %v", err)
	}
}

func TestTUIAddSaves(t *testing.T) {
	isolatedDataDir(t)

	m, err := runTUI(t, "acode\rs", 0)
	if err != nil {
		t.Fatalf("tui add flow: %v", err)
	}
	if len(m.apps) != 1 || m.apps[0].Name != "code" {
		t.Errorf("final model apps = %v, want [{code false}]", m.apps)
	}

	apps, err := LoadApps()
	if err != nil {
		t.Fatalf("LoadApps() error = %v", err)
	}
	if len(apps) != 1 || apps[0].Name != "code" {
		t.Errorf("LoadApps() = %v, want [{code false}]", apps)
	}
}

func TestTUIAddDuplicateShowsError(t *testing.T) {
	isolatedDataDir(t)

	if err := SaveApps([]App{{Name: "rider"}}); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	// "arider\r" leaves the model in add mode with the duplicate error set;
	// there is no quit key in that mode, so settle then quit externally.
	m, err := runTUI(t, "arider\r", 200*time.Millisecond)
	if err != nil {
		t.Fatalf("tui duplicate add: %v", err)
	}
	if m.err == nil {
		t.Fatal("tui duplicate add: err = nil, want duplicate error")
	}
	if !strings.Contains(m.err.Error(), "already exists") {
		t.Errorf("tui duplicate add err = %v, want contains %q", m.err, "already exists")
	}

	apps, _ := LoadApps()
	if len(apps) != 1 || apps[0].Name != "rider" {
		t.Errorf("LoadApps() after failed add = %v, want [{rider false}]", apps)
	}
}

func TestTUIDeleteSaves(t *testing.T) {
	isolatedDataDir(t)

	if err := SaveApps([]App{{Name: "rider"}, {Name: "code"}}); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	if _, err := runTUI(t, "ds", 0); err != nil {
		t.Fatalf("tui delete flow: %v", err)
	}

	apps, _ := LoadApps()
	if len(apps) != 1 || apps[0].Name != "code" {
		t.Errorf("LoadApps() after delete = %v, want [{code false}]", apps)
	}
}

func TestTUIToggleTerminalSaves(t *testing.T) {
	isolatedDataDir(t)

	if err := SaveApps([]App{{Name: "rider"}}); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	if _, err := runTUI(t, "ts", 0); err != nil {
		t.Fatalf("tui toggle flow: %v", err)
	}

	apps, _ := LoadApps()
	if len(apps) != 1 || apps[0].Name != "rider" || !apps[0].IsTerminal {
		t.Errorf("LoadApps() after toggle = %v, want [{rider true}]", apps)
	}
}

func TestTUIQuitWithoutSaving(t *testing.T) {
	isolatedDataDir(t)

	seed := []App{{Name: "rider"}, {Name: "code"}}
	if err := SaveApps(seed); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	// d = delete (dirty), q = confirm dialog, q = quit without saving.
	if _, err := runTUI(t, "dqq", 0); err != nil {
		t.Fatalf("tui delete-then-quit: %v", err)
	}

	apps, _ := LoadApps()
	if len(apps) != 2 {
		t.Errorf("LoadApps() after quit = %v, want unchanged %v", apps, seed)
	}
}

func TestTUIQuitConfirmSaves(t *testing.T) {
	isolatedDataDir(t)

	if err := SaveApps([]App{{Name: "rider"}}); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	// d = delete (dirty), q = confirm dialog, s = save and quit.
	if _, err := runTUI(t, "dqs", 0); err != nil {
		t.Fatalf("tui delete-then-confirm-save: %v", err)
	}

	apps, _ := LoadApps()
	if len(apps) != 0 {
		t.Errorf("LoadApps() after confirm-save = %v, want empty", apps)
	}
}

func TestTUIQuitCleanSkipsConfirm(t *testing.T) {
	isolatedDataDir(t)

	seed := []App{{Name: "rider"}, {Name: "code"}}
	if err := SaveApps(seed); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	// q on a clean list quits immediately, no confirm dialog.
	if _, err := runTUI(t, "q", 0); err != nil {
		t.Fatalf("tui clean quit: %v", err)
	}

	apps, _ := LoadApps()
	if len(apps) != 2 {
		t.Errorf("LoadApps() after clean quit = %v, want unchanged %v", apps, seed)
	}
}

func TestTUIQuitConfirmBack(t *testing.T) {
	isolatedDataDir(t)

	seed := []App{{Name: "rider"}, {Name: "code"}}
	if err := SaveApps(seed); err != nil {
		t.Fatalf("SaveApps() error = %v", err)
	}

	// d = delete (dirty), q = confirm, esc = back to browse, q = confirm again,
	// q = quit without saving. esc is written alone because a lone ESC byte is
	// only resolved after ultraviolet's 50ms escape timeout.
	_, err := runTUIStaged(t,
		[]tuiStage{
			{script: "dq"},
			{script: "\x1b", wait: 150 * time.Millisecond},
			{script: "qq"},
		},
		0,
	)
	if err != nil {
		t.Fatalf("tui confirm-back flow: %v", err)
	}

	apps, _ := LoadApps()
	if len(apps) != 2 {
		t.Errorf("LoadApps() after confirm-back-quit = %v, want unchanged %v", apps, seed)
	}
}

func TestTUIQuitEmpty(t *testing.T) {
	isolatedDataDir(t)

	if _, err := runTUI(t, "q", 0); err != nil {
		t.Fatalf("tui quit on empty list: %v", err)
	}
}
