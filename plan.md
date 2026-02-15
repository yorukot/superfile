## Comparison: superfile `VideoGenerator` vs yazi video previewer

### Where superfile spends time
**Code**: `src/pkg/file_preview/thumbnail_generator.go`

- superfile renders video previews by **spawning `ffmpeg` to generate a JPEG thumbnail**:
  ```go
  ffmpeg -v warning -an -sn -dn -t 180 -hwaccel auto -skip_frame nokey -i <input> \
         -vf thumbnail -frames:v 1 -f image2 -fs 104857600 -y <tmp>.jpg
  ```
- Thumbnail cache is **in-memory only**, and the thumbnails are stored in a **fresh temp dir each launch**:
  - `NewThumbnailGenerator()` -> `os.MkdirTemp("", "superfiles-*")`
  - `tempFilesCache map[string]string` only survives for the current process
  - `CleanUp()` deletes the directory at exit
- Thumbnail generation is currently **synchronous on the render path**:
  - `src/internal/ui/preview/render.go` calls `GetThumbnailOrGenerate(itemPath)` directly inside `RenderWithPath`.
  - There is no cancellation/abort if the user moves selection quickly.

**Implications**
- If users scroll quickly over many videos, superfile can spawn **many `ffmpeg` processes sequentially** and wait for each.
- On every app restart, superfile will regenerate thumbnails again.
- `-vf thumbnail` makes ffmpeg **decode frames to decide “best thumbnail”**. Even with `-skip_frame nokey`, it can still do meaningful work before it outputs a frame.
- superfile does **not explicitly scale** thumbnails down during generation. The scale happens later during terminal rendering (image->ANSI/kitty), which means extra I/O + potentially larger JPEGs than needed.

### Why yazi feels fast
**Code**: `bin/yazi/yazi-plugin/preset/plugins/video.lua` (+ runtime in `bin/yazi/yazi-core/src/tab/preview.rs`)

- yazi’s video previewer is a **plugin** that:
  1. Computes a **stable cache key** based on file identity + `skip` (seek position)
     - `ya.file_cache(job)` -> hashes file + skip into `YAZI.preview.cache_dir` (persistent) (`bin/yazi/yazi-plugin/src/utils/cache.rs`).
  2. If cache file exists and is non-empty, it **reuses it instantly**.
  3. Otherwise it runs `ffprobe` to get duration and decide a seek point, then runs `ffmpeg` with **explicit scaling**:
     ```lua
     ffmpeg -v warning -hwaccel auto -threads 1 -an -sn -dn \
           [-ss <time>] -skip_frame nokey -i <input> \
           -vframes 1 -q:v <quality> \
           -vf scale='min(maxW,iw)':'min(maxH,ih)':force_original_aspect_ratio=decrease:flags=fast_bilinear \
           -f image2 -y <cache>
     ```
  4. It uses `skip` to support **seeking**, and clamps bounds so it doesn’t keep trying beyond available keyframes.
- yazi’s core preview system has **cancellation tokens**:
  - `Preview::go()` aborts previous preview tasks when selection changes (`previewer_ct.cancel()`), so it doesn’t waste time finishing work for a file you already left.
- Cache directory is persistent and ensured to exist: `bin/yazi/yazi-config/src/preview/preview.rs`.

**Implications**
- yazi often does *zero* work after the first time because it hits the persistent cache.
- When the user moves quickly, yazi cancels outdated work aggressively.
- yazi generates thumbnails **already scaled to the preview area**, so less data is written/read/decoded.

## Key performance-critical differences
1. **Persistent cache**
   - superfile: per-run temp dir + in-memory map only
   - yazi: stable cache path in a persistent cache dir

2. **Cancellation / debounce**
   - superfile: thumbnail generation is blocking inside render; no cancel on selection change
   - yazi: cancels previous preview task; also has `image_delay` to debounce image showing

3. **Frame selection strategy**
   - superfile: `-vf thumbnail` (can be expensive)
   - yazi: `ffprobe duration` + `-ss <time>` then grab 1 keyframe

4. **Scaling early**
   - superfile: no `scale` filter in `ffmpeg` stage
   - yazi: always scales to max preview WxH

5. **Threads/CPU predictability**
   - yazi forces `-threads 1` (less contention, especially when rapidly spawning tasks)
   - superfile leaves default ffmpeg threading (can spike CPU and stall UI on busy systems)

## Actionable optimizations for superfile (prioritized)

### Quick wins (very likely to help immediately)
1. **Switch to persistent thumbnail cache**
   - Instead of `os.MkdirTemp`, write to `~/.cache/superfile/thumbnails` (or OS-appropriate cache dir).
   - Use a stable key (hash of full path + mtime/size, maybe inode) so thumbnails survive restarts.
   - Must make sure that we don't fill up user's disk space. Keep the size limited to say 100 MB or a user configured value

2. **Add scaling in `VideoGenerator` ffmpeg command**
   - Mimic yazi:
     - `-vf scale='min(W,iw)':'min(H,ih)':force_original_aspect_ratio=decrease:flags=fast_bilinear`
   - W/H can be based on preview panel size (pass dimensions into generator).

3. **Use seek + vframes instead of `-vf thumbnail`**
   - Example approach:
     - run `ffprobe` once to get duration
     - `-ss <duration*0.05>` (or 0 for attached_pic) + `-vframes 1`
   - This is usually much faster than letting `thumbnail` pick.

### Medium changes
4. **Make thumbnail generation async + cancelable**
   - Don’t generate inside `RenderWithPath`.
   - Kick off generation in a goroutine, render a “Loading…” placeholder, then update preview when ready.
   - Keep a `context.CancelFunc` per currently-previewed file and cancel when selection changes.

5. **Debounce requests when scrolling**
   - Wait e.g. 50–150ms after selection changes before spawning `ffmpeg` (similar to yazi’s `image_delay`).

### Deeper improvements
6. **Pre-generate / prefetch**
   - When hovering a video, pre-generate thumbnails for adjacent items in the list (yazi advertises “precache images and videos”).

7. **Avoid spawning ffmpeg for embedded cover art**
   - For formats with attached_pic, extract the cover directly (`-map disp:attached_pic`) like yazi.

## Next step I suggest
If you want, I can implement the first two quick wins in superfile:
- persistent cache dir + stable cache key
- `ffmpeg` scaling and `-ss` seek strategy

I’ll need one decision from you: should superfile’s cache live under XDG (`$XDG_CACHE_HOME/superfile`) or under the existing superfile config dir (`~/.config/superfile`)?