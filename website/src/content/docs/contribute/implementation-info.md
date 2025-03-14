---
title: Implementation info
description: A collection of general information regarding how various things work
head:
  - tag: title
    content: Implementation info | superfile
---

# Implmentation info
The purpose of this document is to provide some implementation details to the reader that are not so obvious from the code and not very straightforward to figure out. 

## How default configuration files are packaged with app
We use golangs `embed.FS` and embed all files in `src/superfile_config/` into our spf binary. In `src/internal/config_function.go`, the function `LoadAllDefaultConfig()` reads these embedded files, and write them to disk / in memory configuratin variables.
