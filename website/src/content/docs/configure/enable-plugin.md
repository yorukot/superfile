---
title: Enable plugin
description: Enable superfile plguins
head:
  - tag: title
    content: Enable plugins | superfile
---

You can enter the following command to set it up;

[Click me to know where is CONFIG_PATH](/configure/config-file-path#config)

```bash
$EDITOR CONFIG_PATH
```

example:
I want to enable metadata plugin

Please make sure you have installed the Requirements of this plugin.

After that edit `config.toml` using your preferred editor:

```
$EDITOR CONFIG_PATH
```

and change:

```diff
- metadata = false
+ metadata = true
```

## Plugin list

[click me to check plugin list](/list/plugin-list)