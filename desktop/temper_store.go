package main

import (
	"os"
	"path/filepath"
	"sync"

	"reasonix/internal/temper/store"
)

// temperCoWorkDir 返回 CoWork 业务 DB 目录(%APPDATA%\Temper\cowork,
// 由 ApplyTemperIdentity 注入的 REASONIX_STATE_HOME 决定)。
func temperCoWorkDir() string {
	if dir := os.Getenv("REASONIX_STATE_HOME"); dir != "" {
		return dir
	}
	return filepath.Join(temperHomeRoot(), "cowork")
}

// temperStore 是进程级懒加载的 CoWork Store。桌面绑定通过它访问
// Project/Work/Evidence/Artifact 等业务数据。并发安全。
type temperStore struct {
	mu    sync.Mutex
	store *store.Store
	err   error
}

// get 打开(或复用)CoWork Store。失败会缓存错误,避免反复尝试。
func (ts *temperStore) get() (*store.Store, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.store != nil {
		return ts.store, nil
	}
	if ts.err != nil {
		return nil, ts.err
	}
	s, err := store.Open(temperCoWorkDir())
	if err != nil {
		ts.err = err
		return nil, err
	}
	ts.store = s
	return s, nil
}

// close 关闭 Store(应用退出时调用)。
func (ts *temperStore) close() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.store != nil {
		_ = ts.store.Close()
		ts.store = nil
	}
}
