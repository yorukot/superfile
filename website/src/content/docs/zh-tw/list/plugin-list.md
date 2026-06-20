---
title: 外掛清單
description: 可用 superfile 外掛的完整清單
head:
  - tag: title
    content: 外掛清單 | superfile
---

superfile 支援多種外掛來擴充功能。以下是可用外掛及其需求的完整清單。

### Metadata

- **描述：** 顯示檔案與目錄更詳細的 metadata

- **需求：** [`exiftool`](https://exiftool.org)

- **設定名稱：** `metadata`

### MD5 Checksum

- **描述：** 在 metadata panel 顯示一般檔案的 MD5 checksum

- **需求：** 無

- **設定名稱：** `enable_md5_checksum`

- **注意：** 計算 checksum 會讀取選取的檔案，檔案很大時可能會比較慢。

### Zoxide

- **描述：** 與 zoxide 整合的智慧目錄跳轉功能。透過可搜尋的視窗快速導覽到常用目錄。

- **需求：** [`zoxide`](https://github.com/ajeetdsouza/zoxide)

- **設定名稱：** `zoxide_support`

- **使用方式：** 按下 `z` 開啟 zoxide 導覽視窗。開始輸入以搜尋目錄，使用方向鍵在結果間移動，並按 Enter 跳轉到目錄。
