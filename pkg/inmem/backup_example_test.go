package inmem

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var _ = describe("InitBackup", func(t *testing.T) {
	It("initializes the backup process successfully", func(t *testing.T) {
		// Create a temporary file for backup
		f, err := os.CreateTemp("", "example-backup-file")
		require.NoError(t, err)
		defer os.Remove(f.Name()) // Clean up

		// Initialize the storage with some data
		storage := New[string, string]()
		storage.Save(context.Background(), "key1", "value1")
		storage.Save(context.Background(), "key2", "value2")

		// Initialize backup
		err = storage.InitBackup(context.Background(), f, 5*time.Second, nil)
		require.NoError(t, err)

		// Check that the backup loop runs
		time.Sleep(10 * time.Millisecond) // Give it some time to start the loop

		// Verify that the backup file contains the expected data
		data, err := os.ReadFile(f.Name())
		require.NoError(t, err)
		require.True(t, strings.Contains(string(data), `"key1":"value1"`))
		require.True(t, strings.Contains(string(data), `"key2":"value2"`))
	})
})

var _ = describe("LoadFromBackup", func(t *testing.T) {
	It("loads data from the backup file into the store", func(t *testing.T) {
		// Create a temporary file for backup
		f, err := os.CreateTemp("", "example-backup-file")
		require.NoError(t, err)
		defer os.Remove(f.Name()) // Clean up

		// Initialize the storage with some data
		storage := New[string, string]()
		storage.Save(context.Background(), "key1", "value1")
		storage.Save(context.Background(), "key2", "value2")

		// Dump the current state of the store to the backup file
		err = storage.dumpBackup()
		require.NoError(t, err)

		// Create a new storage instance
		newStorage := New[string, string]()

		// Load the backup into the new storage
		err = newStorage.LoadFromBackup()
		require.NoError(t, err)

		// Verify that the new storage has the same data as the original one
		val1, found1 := newStorage.Get(context.Background(), "key1")
		require.True(t, found1)
		require.Equal(t, "value1", val1)

		val2, found2 := newStorage.Get(context.Background(), "key2")
		require.True(t, found2)
		require.Equal(t, "value2", val2)
	})
})

var _ = describe("shutBackup", func(t *testing.T) {
	It("closes the backup file gracefully", func(t *testing.T) {
		// Create a mock backup file using bytes.Buffer
		mockBackupFile := &mockFile{}

		// Initialize the storage with a mock backup file
		storage := New[string, string]()
		storage.backupFile = mockBackupFile

		// Call shutBackup
		storage.shutBackup()

		// Verify that the Close method was called on the mock backup file
		require.True(t, mockBackupFile.closed)
	})
})

// mockFile is used to simulate a file for testing purposes.
type mockFile struct {
	closed bool
}

func (m *mockFile) Read(p []byte) (int, error) {
	panic("mockFile does not implement Read")
}

func (m *mockFile) Write(p []byte) (int, error) {
	panic("mockFile does not implement Write")
}

func (m *mockFile) Truncate(_ int64) error {
	panic("mockFile does not implement Truncate")
}

func (m *mockFile) Seek(_ int64, _ int) (int64, error) {
	panic("mockFile does not implement Seek")
}

func (m *mockFile) Close() error {
	m.closed = true
	return nil
}

// describe is a helper function to structure tests similar to BDD style.
func describe(description string, spec func(t *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Run(description, spec)
	}
}

// It is a helper function to define individual test cases within a describe block.
func It(text string, body func(t *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Run(text, body)
	}
}
