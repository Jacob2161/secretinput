package secretinput

import (
	"bytes"
	"os"
	"strconv"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

// openTerminal opens a pseudo-terminal and returns its controlling side,
// which plays the person at the keyboard, and its terminal side, which the
// program reads.
func openTerminal(t *testing.T) (keyboard *os.File, terminal *os.File) {
	t.Helper()
	keyboard, err := os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Skipf("no pseudo-terminals here: %v", err)
	}
	t.Cleanup(func() { _ = keyboard.Close() })
	descriptor := int(keyboard.Fd())
	if err := unix.IoctlSetPointerInt(descriptor, unix.TIOCSPTLCK, 0); err != nil {
		t.Fatalf("unlocking the pseudo-terminal: %v", err)
	}
	number, err := unix.IoctlGetInt(descriptor, unix.TIOCGPTN)
	if err != nil {
		t.Fatalf("naming the pseudo-terminal: %v", err)
	}
	path := "/dev/pts/" + strconv.Itoa(number)
	terminal, err = os.OpenFile(path, os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	t.Cleanup(func() { _ = terminal.Close() })
	return keyboard, terminal
}

// echoing reports whether terminal echoes what is typed.
func echoing(t *testing.T, terminal *os.File) bool {
	t.Helper()
	settings, err := unix.IoctlGetTermios(int(terminal.Fd()), unix.TCGETS)
	if err != nil {
		t.Fatalf("reading the terminal settings: %v", err)
	}
	return settings.Lflag&unix.ECHO != 0
}

// TestPromptOnTerminal types a secret into a real pseudo-terminal and checks
// that Prompt returns it, that nothing typed was echoed back, that the
// prompt and its line end went to the writer, and that echo is back on
// afterwards.
func TestPromptOnTerminal(t *testing.T) {
	keyboard, terminal := openTerminal(t)
	if !echoing(t, terminal) {
		t.Fatal("a new pseudo-terminal does not echo, want echo on")
	}
	type promptAnswer struct {
		secret string
		err    error
	}
	answers := make(chan promptAnswer, 1)
	var prompt bytes.Buffer
	go func() {
		secret, err := Prompt(terminal, &prompt, "Password: ")
		answers <- promptAnswer{secret: secret, err: err}
	}()
	// Type only once Prompt has turned echo off, as a person would after
	// seeing the prompt; typing earlier would be echoed by design.
	deadline := time.Now().Add(10 * time.Second)
	for echoing(t, terminal) {
		if time.Now().After(deadline) {
			t.Fatal("Prompt never turned echo off")
		}
		time.Sleep(time.Millisecond)
	}
	if _, err := keyboard.WriteString("hunter2\n"); err != nil {
		t.Fatal(err)
	}
	answer := <-answers
	if answer.err != nil || answer.secret != "hunter2" {
		t.Fatalf("Prompt = %q, %v; want hunter2", answer.secret, answer.err)
	}
	if got := prompt.String(); got != "Password: \n" {
		t.Errorf("the writer got %q, want %q", got, "Password: \n")
	}
	if !echoing(t, terminal) {
		t.Error("echo is still off after Prompt returned, want it restored")
	}
	// Everything the terminal displayed comes back on the keyboard side; a
	// marker written after Prompt ends what to read.
	if _, err := terminal.WriteString("marker\n"); err != nil {
		t.Fatal(err)
	}
	var displayed bytes.Buffer
	reading := make([]byte, 256)
	for !bytes.Contains(displayed.Bytes(), []byte("marker")) {
		n, err := keyboard.Read(reading)
		if err != nil {
			t.Fatalf("reading what the terminal displayed: %v", err)
		}
		displayed.Write(reading[:n])
	}
	if bytes.Contains(displayed.Bytes(), []byte("hunter2")) {
		t.Fatalf("the terminal displayed %q, want the secret left out", displayed.String())
	}
}
