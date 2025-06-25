---
title: Enable Plugin
description: How to enable and configure superfile plugins
head:
  - tag: title
    content: Enable Plugins | superfile
---

Plugins extend superfile's functionality by integrating with external tools. This guide shows you how to enable and configure plugins.

## Prerequisites

Before enabling any plugin, ensure you have:

1. **Installed the required dependencies** for the specific plugin
2. **Located your config file** - see [config file path guide](/configure/config-file-path#config)

## How to Enable Plugins

### Step 1: Install Required Dependencies

Each plugin has specific requirements. Check the [plugin list](/list/plugin-list) for the dependencies needed for your desired plugin.

### Step 2: Edit Configuration File

Open your `config.toml` file:

```bash
$EDITOR CONFIG_PATH
```

### Step 3: Enable the Plugin

Find the plugin section in your config and change its value from `false` to `true`:

```diff
[plugins]
- metadata = false
+ metadata = true
```

### Example: Enabling Metadata Plugin

1. **Install exiftool** (required for metadata plugin)
2. **Edit your config file:**
   ```bash
   $EDITOR CONFIG_PATH
   ```
3. **Enable the plugin:**
   ```toml
   metadata = true
   ```

## Configuration Format

```toml
metadata = false
zoxide = false
```

Set any plugin to `true` to enable it, or `false` to disable it.

## Available Plugins

For a complete list of available plugins and their requirements, see the [plugin list](/list/plugin-list).

## Troubleshooting

If a plugin isn't working after enabling it:

1. **Verify dependencies** - Make sure all required tools are installed and accessible in your PATH
2. **Restart superfile** - Changes require restarting the application
3. **Check configuration** - Ensure the plugin name is spelled correctly in your config file