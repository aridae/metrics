package inmem

import (
	"context"
	"testing"
	"time"

	"github.com/aridae/go-metrics-store/pkg/inmem/_mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestNew проверяет создание нового хранилища.
func TestNew(t *testing.T) {
	s := New[int, string]()
	require.NotNil(t, s)
	require.Equal(t, make(map[int]string), s.store)
}

// TestInitBackup проверяет инициализацию резервного копирования.
func TestInitBackup(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockFile := _mock.NewMockfile(ctrl)

	s := New[int, string]()

	err := s.InitBackup(ctx, mockFile, 1*time.Second, nil)

	require.NoError(t, err)
	require.Equal(t, mockFile, s.backupFile)
	require.Equal(t, 1*time.Second, s.backupInterval)
}

// TestDumpBackup проверяет процесс создания резервной копии.
func TestDumpBackup(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockFile := _mock.NewMockfile(ctrl)
	mockEncoder := _mock.NewMockencoder(ctrl)

	s := New[int, string]()
	s.store = map[int]string{1: "value"}
	s.provideFileEncoder = func(file) encoder { return mockEncoder }

	mockFile.EXPECT().Truncate(int64(0)).Return(nil)
	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), nil)
	mockEncoder.EXPECT().Encode(s.store).Return(nil)

	err := s.InitBackup(ctx, mockFile, 1*time.Second, nil)
	require.NoError(t, err)

	err = s.dumpBackup()
	require.NoError(t, err)
}

// TestLoadFromBackup проверяет загрузку данных из резервной копии.
func TestLoadFromBackup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := _mock.NewMockfile(ctrl)
	mockDecoder := _mock.NewMockdecoder(ctrl)

	s := New[int, string]()
	s.backupFile = mockFile
	s.providerFileDecoder = func(file) decoder { return mockDecoder }

	mockDecoder.EXPECT().Decode(gomock.Any()).Return(nil)

	err := s.LoadFromBackup()
	require.NoError(t, err)
}
