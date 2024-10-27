package inmem

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"io"
	"os"
	"time"
)

func (mem *MemTimeseriesStorage) InitBackup(
	ctx context.Context,
	backupFilepath string,
	backupInterval time.Duration,
	registerTypes map[string]any,
) error {
	backupFile, err := os.OpenFile(backupFilepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file for backup %s: %w", backupFilepath, err)
	}

	mem.backupFile = backupFile
	mem.backupInterval = backupInterval

	// register types for gob serialization
	for name, val := range registerTypes {
		gob.RegisterName(name, val)
	}

	go mem.runBackupLoop(ctx)

	return nil
}

func (mem *MemTimeseriesStorage) runBackupLoop(ctx context.Context) {
	ticker := time.NewTicker(mem.backupInterval)
	defer func() {
		mem.shutBackup()
	}()

	for {
		select {
		case <-ticker.C:
			if err := mem.dumpBackup(); err != nil {
				logger.Obtain().Errorf("[timeseriesstorage.MemTimeseriesStorage.runBackupLoop][CRITICAL] failed to dump data to backup file: %v", err)
			}
		case <-ctx.Done():
			logger.Obtain().Info("stopping backup service downstreams...")
			return
		}
	}
}

func (mem *MemTimeseriesStorage) shutBackup() {
	mem.fileMu.Lock()
	defer mem.fileMu.Unlock()

	logger.Obtain().Info("closing backup file...")
	err := mem.backupFile.Close()
	if err != nil {
		logger.Obtain().Errorf("failed to close backup file: %v", err)
	}

	logger.Obtain().Info("backup is shut")
}

func (mem *MemTimeseriesStorage) dumpBackup() error {
	mem.fileMu.Lock()
	mem.storeMu.RLock()
	defer mem.fileMu.Unlock()
	defer mem.storeMu.RUnlock()

	// NOTE: а есть способ сделать это более элегантно? Рыдаю ToT
	mem.backupFile.Truncate(0) //nolint:errcheck
	mem.backupFile.Seek(0, 0)  //nolint:errcheck

	return gob.NewEncoder(mem.backupFile).Encode(mem.store)
}

func (mem *MemTimeseriesStorage) LoadFromBackup() error {
	if mem.backupFile == nil {
		return fmt.Errorf("no backup file found, make sure to call InitBackup() method to init backing options before loading from backup file")
	}

	mem.fileMu.RLock()
	mem.storeMu.Lock()
	defer mem.fileMu.RUnlock()
	defer mem.storeMu.Unlock()

	newStore := mem.store
	err := gob.NewDecoder(mem.backupFile).Decode(&newStore)
	if err == io.EOF {
		logger.Obtain().Info("backup file is empty, nothing to load")
		return nil
	}

	_ = newStore

	return err
}
