---
title: Image Preview
description: Learn how image preview works in superfile and how terminal compatibility is determined.
head:
  - tag: title
    content: Image Preview | superfile
---

This tutorial will teach you how to use superfile’s image preview feature step by step.

## What is Image Preview?

superfile supports image previews directly in your terminal using several display protocols. When supported, images can be shown inline without any external viewer.

---

## Terminal Compatibility

superfile automatically detects your terminal using the `$TERM` and `$TERM_PROGRAM` environment variables. We support rendering on the following terminals:

| Terminal              | Protocol         | Image Preview Support |
|-----------------------|------------------|------------------------|
| **kitty**             | Kitty protocol   | ✅                     |
| **WezTerm**           | Kitty protocol   | ✅                     |
| **Ghostty**           | Kitty protocol   | ✅                     |
| **iTerm2**            | Inline images    | ❌                     |
| **Konsole**           | Inline images    | ❌                     |
| **VSCode**            | Inline images    | ❌                     |
| **Tabby**             | Inline images    | ❌                     |
| **Hyper**             | Inline images    | ❌                     |
| **Mintty**            | Inline images    | ❌                     |
| **foot**              | Sixel graphics   | ❌                     |
| **Black Box**         | Sixel graphics   | ❌                     |

> ✅ means full support for inline image preview using Kitty protocol  
> ❌ means image preview is currently not supported

---

## Supported Protocols

superfile supports the following rendering protocols and will automatically choose the best one based on your terminal:

| Protocol Name     | Description                                                                                   | Status      |
|-------------------|-----------------------------------------------------------------------------------------------|-------------|
| **Kitty protocol** | Most capable, pixel-accurate rendering with transparency and scaling support.                | ✅ Preferred|
| **Sixel**          | Old standard used in DEC terminals and some modern ones like foot.                           | ❌          |
| **iTerm2 inline**  | iTerm2’s proprietary image format, used in Tabby, Hyper, etc.                                | ❌          |
| **ANSI**           | Fallback text rendering using ANSI blocks or metadata only.                                  | ✅ Always   |

---

## Terminal Detection and Pixel Size

superfile detects terminal capabilities by inspecting:

- `$TERM`
- `$TERM_PROGRAM`

These variables help us decide whether advanced rendering might be possible. However, real support is confirmed at runtime using terminal queries.

To scale images correctly, superfile sends the following escape code:

```
\x1b[16t
```

This sequence queries the terminal for the size of each **cell in pixels**. superfile uses the result to:

- Maintain correct image aspect ratio
- Avoid distortions in previews
- Adapt to terminal resizes

If your terminal does not support `\x1b[16t`, we fallback to default assumptions like `10×20 px per cell`.

## Graceful Fallback to ANSI

When advanced image preview isn't supported (for example, when the terminal doesn't support the Kitty protocol), superfile gracefully falls back to an ANSI-based preview using color-coded blocks.

This ensures a consistent and reliable experience across all terminal environments.