// Package main — Temper 身份与数据隔离。
//
// Temper 构建在 Reasonix Runtime 之上,但必须拥有自己的产品身份与数据
// 命名空间。本文件在 Reasonix Boot / Config load 之前注入 Temper 的
// Runtime Home(REASONIX_HOME / REASONIX_STATE_HOME / REASONIX_CACHE_HOME),
// 确保 Temper 不读取/写入正式 Reasonix 的 Provider、Session、Memory、
// Plugin、Cache。
//
// 用户可见身份:
//
//	产品     Temper
//	版本     0.3.0-dev
//	标语     Shape intent. Ship work.
//	可执行   Temper.exe
//
// 数据位置:
//
//	%APPDATA%\Temper\runtime\    REASONIX_HOME(运行时配置/状态)
//	%APPDATA%\Temper\cowork\     REASONIX_STATE_HOME(CoWork 业务数据)
//	%LOCALAPPDATA%\Temper\cache\ REASONIX_CACHE_HOME(缓存/派生数据)
//
// 网络审计(PHASE B):
//
//	Provider network       RETAIN
//	User Web tools         RETAIN
//	User MCP               RETAIN
//	Reasonix telemetry     DISABLE(默认)
//	Reasonix crash upload  DISABLE(默认)
//	Reasonix updater       REPLACE/DISABLE(Temper 自管更新通道)
//	Reasonix release URLs  REPLACE(Temper 发布)
//
// 实现说明:这是 Temper 身份层(adapter),不修改 Reasonix-owned source。
// 对用户可见的字符串(窗口标题、托盘、主题作者等)在各自的 UI 面替换。
package main

import (
	"os"
	"path/filepath"
)

const (
	// TemperDisplayName 是用户可见的产品名。
	TemperDisplayName = "Temper"
	// TemperVersion 是当前开发版本号。
	TemperVersion = "0.3.0-dev"
	// TemperTagline 是产品标语。
	TemperTagline = "Shape intent. Ship work."
	// TemperExecutable 是 Windows 可执行文件名(不含扩展名)。
	TemperExecutable = "Temper"
)

// temperHomeRoot 返回 Temper 的运行时根目录。
func temperHomeRoot() string {
	if dir := os.Getenv("TEMPER_HOME"); dir != "" {
		return dir
	}
	if appData, err := os.UserConfigDir(); err == nil && appData != "" {
		return filepath.Join(appData, "Temper")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "AppData", "Roaming", "Temper")
	}
	return ""
}

// temperCacheRoot 返回 Temper 的缓存根目录。
func temperCacheRoot() string {
	if dir := os.Getenv("TEMPER_CACHE"); dir != "" {
		return dir
	}
	if cache, err := os.UserCacheDir(); err == nil && cache != "" {
		return filepath.Join(cache, "Temper")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "AppData", "Local", "Temper")
	}
	return ""
}

// ApplyTemperIdentity 在 Reasonix Boot / Config load 之前调用,注入 Temper
// 的 Runtime Home 隔离。幂等:已显式设置的环境变量不会被覆盖。
func ApplyTemperIdentity() {
	// 只注入未设置的值:用户/CI 显式设置的 REASONIX_* 优先。
	setDefaultEnv("REASONIX_HOME", filepath.Join(temperHomeRoot(), "runtime"))
	setDefaultEnv("REASONIX_STATE_HOME", filepath.Join(temperHomeRoot(), "cowork"))
	setDefaultEnv("REASONIX_CACHE_HOME", filepath.Join(temperCacheRoot(), "cache"))

	// Temper 默认关闭 Reasonix 远程产品遥测与崩溃上报。
	setDefaultEnv("REASONIX_TELEMETRY", "0")
	setDefaultEnv("DO_NOT_TRACK", "1")
	// Reasonix 崩溃上传默认关闭。
	setDefaultEnv("REASONIX_CRASH_UPLOAD", "0")
	// Reasonix 自动更新器由 Temper 自管,默认禁用其远端 manifest 检查。
	setDefaultEnv("REASONIX_UPDATE_DISABLE", "1")
}

// TemperDefaultTelemetryOff 返回 Temper 的产品遥测默认策略:默认关闭。
// 用户可在 Settings 中显式开启。
const TemperDefaultTelemetryOff = true

// setDefaultEnv 仅在变量未设置时写入环境。
func setDefaultEnv(key, value string) {
	if _, ok := os.LookupEnv(key); !ok {
		_ = os.Setenv(key, value)
	}
}
