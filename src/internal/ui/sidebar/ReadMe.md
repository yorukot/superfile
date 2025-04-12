# sidebar package
This is for the sidebar UI, and for fetching and updating sidebar directories

# To-dos
- Add missing unit tests
- Separate out implementation of file I/O operations. (Disk listing, Reading and Updating pinned.json)
  This package should only be concerned with UI/UX.
- Implementing a proper state transitioning for the sidebar's different modes (normal, search, rename)
- Some methods could be made more pure by reducing side effects

# Coverage
# Coverage

```bash
cd /path/to/ui/prompt
go test -cover
```
Current coverage is 74.0%.
Current coverage