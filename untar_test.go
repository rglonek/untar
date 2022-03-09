package untar

import (
	"os"
	"testing"
)

func TestUntarSmall(t *testing.T) {
	err := UntarFile("testsmall.tgz", "./")
	os.RemoveAll("./testdir")
	if err != nil {
		t.Logf("failed untar: %s", err)
		t.FailNow()
	}
}

func TestUntarBz(t *testing.T) {
	err := UntarFile("testbz.tgz", "./")
	os.RemoveAll("./testbz")
	if err != nil {
		t.Logf("failed untar: %s", err)
		t.FailNow()
	}
}

func TestUntar(t *testing.T) {
	err := UntarFile("testtar.tar", "./")
	os.RemoveAll("./testtar")
	if err != nil {
		t.Logf("failed untar: %s", err)
		t.FailNow()
	}
}
