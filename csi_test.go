package terminal

import (
	"testing"
)

func TestCSIParse(t *testing.T) {
	var csi csiEscape
	csi.reset()
	csi.buf = []byte("s")
	csi.parse()
	if csi.mode != 's' || csi.arg(0, 17) != 17 || len(csi.args) != 0 {
		t.Fatal("CSI parse mismatch")
	}

	csi.reset()
	csi.buf = []byte("31T")
	csi.parse()
	if csi.mode != 'T' || csi.arg(0, 0) != 31 || len(csi.args) != 1 {
		t.Fatal("CSI parse mismatch")
	}

	csi.reset()
	csi.buf = []byte("48;2f")
	csi.parse()
	if csi.mode != 'f' || csi.arg(0, 0) != 48 || csi.arg(1, 0) != 2 || len(csi.args) != 2 {
		t.Fatal("CSI parse mismatch")
	}

	csi.reset()
	csi.buf = []byte("?25l")
	csi.parse()
	if csi.mode != 'l' || csi.arg(0, 0) != 25 || csi.priv != true || len(csi.args) != 1 {
		t.Fatal("CSI parse mismatch")
	}
}

// Test that ESC[=u is not treated as ESC[u
func TestCSIEqualsPrefix(t *testing.T) {
	var st State
	term, err := Create(&st, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Put the cursor somewhere we can detect a restore.
	st.moveTo(5, 5)
	st.saveCursor()
	st.moveTo(10, 10)

	// ESC[=u must not be treated as ESC[u.
	_, err = term.Write([]byte("\033[=u"))
	if err != nil {
		t.Fatal(err)
	}

	x, y := st.Cursor()
	if x != 10 || y != 10 {
		t.Fatalf("ESC[=u moved cursor to (%d, %d), expected (10, 10)", x, y)
	}
}

// Test the difference in behavior between ESC[u and ESC[=u
func TestCSIRestoreCursor(t *testing.T) {
	var st State
	term, err := Create(&st, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Save the cursor at (5, 5), then move it elsewhere.
	st.moveTo(5, 5)
	st.saveCursor()
	st.moveTo(10, 10)

	// ESC[u should restore the saved cursor position.
	_, err = term.Write([]byte("\033[u"))
	if err != nil {
		t.Fatal(err)
	}

	x, y := st.Cursor()
	if x != 5 || y != 5 {
		t.Fatalf("ESC[u moved cursor to (%d, %d), expected (5, 5)", x, y)
	}
}
