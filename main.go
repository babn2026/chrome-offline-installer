package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/babn2026/chrome-offline-installer/do"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	// ========== 1. 加载本地数据 ==========
	data := do.LoadSavedData()

	// ========== 2. 拉取远程最新数据 ==========
	do.FetchAndUpdateData(data)

	// ========== 3. 保存数据到本地 ==========
	if err := do.SaveData(data); err != nil {
		slog.Error(fmt.Sprintf("save data error: %v", err))
	}

	// ========== 4. 保存 Markdown 文档 ==========
	if err := do.SaveMarkdown(data); err != nil {
		slog.Error(fmt.Sprintf("save markdown error: %v", err))
	}

	// ========== 5. 检查是否需要下载 ==========
	// 先取 x64 版本，若不存在或没有下载链接则直接退出
	x64Ver, ok := data["win_stable_x64"]
	if !ok {
		slog.Error("win_stable_x64 not found in data")
		os.Exit(1)
	}
	if len(x64Ver.Urls) == 0 {
		slog.Error("win_stable_x64 has no download URLs")
		os.Exit(1)
	}

	// 以 x64 的版本号判断是否需要下载（x64 与 arm64 通常同步发布）
	if !do.CheckUpdate(x64Ver.Version) {
		slog.Info("no need to download")
		return
	}

	// ========== 6. 下载 x64 ==========
	if err := do.Download("x64", x64Ver.Urls[0]); err != nil {
		slog.Error(fmt.Sprintf("download x64 error: %v", err))
	}

	// ========== 7. 下载 arm64 ==========
	arm64Ver, ok := data["win_stable_arm64"]
	if !ok {
		slog.Error("win_stable_arm64 not found in data")
		os.Exit(1)
	}
	if len(arm64Ver.Urls) == 0 {
		slog.Error("win_stable_arm64 has no download URLs")
		os.Exit(1)
	}

	if err := do.Download("arm64", arm64Ver.Urls[0]); err != nil {
		slog.Error(fmt.Sprintf("download arm64 error: %v", err))
	}
}
