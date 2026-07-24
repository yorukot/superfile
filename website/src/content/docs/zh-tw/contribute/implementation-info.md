---
title: 實作資訊
description: 關於各項機制如何運作的一般資訊集合
head:
  - tag: title
    content: 實作資訊 | superfile
---

# 實作資訊

這份文件的目的是提供一些從程式碼中不太明顯、也不容易直接推斷出的實作細節。

## 預設設定檔如何與應用程式一起封裝

我們使用 Go 的 `embed.FS`，並將 `src/superfile_config/` 中的所有檔案嵌入到 spf binary。於 `src/internal/common/load_config.go` 中，`LoadAllDefaultConfig()` 函式會讀取這些嵌入檔案，並將它們寫入磁碟或記憶體中的設定變數。
