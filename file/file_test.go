package file

import (
	"testing"
	"os"
	"io"
	"io/ioutil"
	"math/rand"
	"crypto/sha256"
	"path/filepath"
	"bytes"
)

const TESTDATA_LEN = 16 * 1024
const TEST_FILENAME = "testfile"

func TestSafeWriteBasic(t *testing.T) {
	dir, err := ioutil.TempDir("", "tmp")
	if err != nil {
		t.Errorf("ioutil.TempDir() failed: %v", err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, TEST_FILENAME)

	data := make([]byte, TESTDATA_LEN)
	_, _ = rand.Read(data)

	// Reference hash
	refHash := sha256.New()
	refHash.Write(data)
	origSum := refHash.Sum(nil)

	err = SafeWrite(path, data)
	if err != nil {
		t.Errorf("SafeWrite() failed: %v", err)
	}

	// Drop written buffer
	data = nil

	f, err := os.Open(filepath.Join(dir, TEST_FILENAME))
	if err != nil {
		t.Errorf("Can't open written file: %v", err)
	}

	// Validate written file contents
	examinedHash := sha256.New()
	_, _ = io.Copy(examinedHash, f)
	examinedSum := examinedHash.Sum(nil)

	if bytes.Compare(origSum, examinedSum) != 0 {
		t.Errorf("File content doesn't match original data: %x != %x",
			origSum, examinedSum)
	}
}
