package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestApplyTemperIdentityInjectsTemperHomes 验证 ApplyTemperIdentity 在
// Reasonix config load 前注入 Temper 的 Runtime Home,且不覆盖用户显式
// 设置的值。
func TestApplyTemperIdentityInjectsTemperHomes(t *testing.T) {
	home := t.TempDir()
	cache := t.TempDir()
	t.Setenv("TEMPER_HOME", home)
	t.Setenv("TEMPER_CACHE", cache)
	// 清除可能影响的环境变量,保证测试确定性。
	unsetAll(t, "REASONIX_HOME", "REASONIX_STATE_HOME", "REASONIX_CACHE_HOME")

	ApplyTemperIdentity()

	if got := os.Getenv("REASONIX_HOME"); got != filepath.Join(home, "runtime") {
		t.Fatalf("REASONIX_HOME = %q, want %q", got, filepath.Join(home, "runtime"))
	}
	if got := os.Getenv("REASONIX_STATE_HOME"); got != filepath.Join(home, "cowork") {
		t.Fatalf("REASONIX_STATE_HOME = %q, want %q", got, filepath.Join(home, "cowork"))
	}
	if got := os.Getenv("REASONIX_CACHE_HOME"); got != filepath.Join(cache, "cache") {
		t.Fatalf("REASONIX_CACHE_HOME = %q, want %q", got, filepath.Join(cache, "cache"))
	}

	// Temper 默认关闭遥测与更新。
	if got := os.Getenv("REASONIX_TELEMETRY"); got != "0" {
		t.Fatalf("REASONIX_TELEMETRY = %q, want 0", got)
	}
	if got := os.Getenv("DO_NOT_TRACK"); got != "1" {
		t.Fatalf("DO_NOT_TRACK = %q, want 1", got)
	}
	if got := os.Getenv("REASONIX_UPDATE_DISABLE"); got != "1" {
		t.Fatalf("REASONIX_UPDATE_DISABLE = %q, want 1", got)
	}
}

// TestApplyTemperIdentityRespectsExplicitEnv 验证已显式设置的 REASONIX_*
// 不被覆盖(用户/CI 显式隔离优先)。
func TestApplyTemperIdentityRespectsExplicitEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TEMPER_HOME", home)
	explicit := filepath.Join(home, "explicit-home")
	t.Setenv("REASONIX_HOME", explicit)

	ApplyTemperIdentity()

	if got := os.Getenv("REASONIX_HOME"); got != explicit {
		t.Fatalf("REASONIX_HOME = %q, want explicit %q (must not be overwritten)", got, explicit)
	}
}

// TestTemperIdentityConstants 验证用户可见身份常量。
func TestTemperIdentityConstants(t *testing.T) {
	if TemperDisplayName != "Temper" {
		t.Fatalf("TemperDisplayName = %q, want Temper", TemperDisplayName)
	}
	if TemperVersion != "0.3.0-dev" {
		t.Fatalf("TemperVersion = %q, want 0.3.0-dev", TemperVersion)
	}
	if !strings.Contains(TemperTagline, "Shape intent") || !strings.Contains(TemperTagline, "Ship work") {
		t.Fatalf("TemperTagline = %q, want 'Shape intent. Ship work.'", TemperTagline)
	}
	if TemperExecutable != "Temper" {
		t.Fatalf("TemperExecutable = %q, want Temper", TemperExecutable)
	}
}

// TestTemperDataIsolationFromReasonix 验证两个不同 Home 下,Temper 的
// REASONIX_HOME 指向 Temper 专属目录,绝不会解析到正式 Reasonix 的
// %APPDATA%\reasonix 或 ~/.reasonix。
func TestTemperDataIsolationFromReasonix(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TEMPER_HOME", home)
	unsetAll(t, "REASONIX_HOME", "REASONIX_STATE_HOME", "REASONIX_CACHE_HOME")

	ApplyTemperIdentity()

	runtimeHome := os.Getenv("REASONIX_HOME")
	if runtimeHome == "" {
		t.Fatal("REASONIX_HOME not set after ApplyTemperIdentity")
	}
	// 必须落在 TEMPER_HOME 之下,绝不落在 reasonix 默认目录。
	abs, err := filepath.Abs(runtimeHome)
	if err != nil {
		t.Fatal(err)
	}
	homeAbs, err := filepath.Abs(home)
	if err != nil {
		t.Fatal(err)
	}
	rel, err := filepath.Rel(homeAbs, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		t.Fatalf("REASONIX_HOME %q must be under TEMPER_HOME %q (rel=%q err=%v)", abs, homeAbs, rel, err)
	}
	// 路径组件中不得出现 reasonix 默认目录名(.reasonix / reasonix)。
	for part := range strings.SplitSeq(strings.ToLower(abs), string(filepath.Separator)) {
		if part == ".reasonix" || part == "reasonix" {
			t.Fatalf("REASONIX_HOME %q must not resolve into a reasonix default directory", abs)
		}
	}
}

func unsetAll(t *testing.T, keys ...string) {
	t.Helper()
	for _, k := range keys {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("unset %s: %v", k, err)
		}
		t.Cleanup(func() { _ = os.Unsetenv(k) })
	}
}
