---
title: 啟用外掛
description: 如何啟用並設定 superfile 外掛
head:
  - tag: title
    content: 啟用外掛 | superfile
---

外掛會透過整合外部工具來擴充 superfile 的功能。本指南會說明如何啟用與設定外掛。

## 前置需求

啟用任何外掛之前，請確認你已經：

1. **安裝特定外掛所需的相依套件**
2. **找到你的設定檔** - 請參考[設定檔路徑指南](/zh-tw/configure/config-file-path#config)

## 如何啟用外掛

### 步驟 1：安裝必要相依套件

每個外掛都有特定需求。請查看[外掛列表](/zh-tw/list/plugin-list)，確認你想使用的外掛需要哪些相依套件。

### 步驟 2：編輯設定檔

開啟你的 `config.toml` 檔案：

```bash
$EDITOR CONFIG_PATH
```

### 步驟 3：啟用外掛

在設定檔中找到 plugin 設定，並將其值從 `false` 改為 `true`：

```diff
- metadata = false
+ metadata = true
```

### 範例：啟用 Metadata 外掛

1. **安裝 exiftool**（metadata 外掛需要）
2. **編輯你的設定檔：**
   ```bash
   $EDITOR CONFIG_PATH
   ```
3. **啟用外掛：**
   ```toml
   metadata = true
   ```

## 設定格式

```toml
metadata = false
enable_md5_checksum = false
zoxide_support = false
```

將任何外掛設為 `true` 即可啟用，或設為 `false` 來停用。

## 可用外掛

如需完整的可用外掛與需求列表，請查看[外掛列表](/zh-tw/list/plugin-list)。

## 疑難排解

如果外掛啟用後無法運作：

1. **確認相依套件** - 確保所有必要工具都已安裝，且可在你的 PATH 中存取
2. **重新啟動 superfile** - 變更需要重新啟動應用程式才會生效
3. **檢查設定** - 確認設定檔中的外掛名稱拼寫正確
