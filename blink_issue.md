# Blink test issue

## Why the old assertion no longer works

The old test asserted `m.textInput.Cursor.Blink` after each blink update. That was valid in Bubbles v1 because:

- `textinput.Model` exposed a public `Cursor cursor.Model` field at `charm.land/bubbles@v1.0.0/textinput/textinput.go:97`.
- That cursor model had a runtime `Blink bool` field at `charm.land/bubbles@v1.0.0/cursor/cursor.go:69`.
- Each blink update toggled that field at `charm.land/bubbles@v1.0.0/cursor/cursor.go:128`.

In Bubbles v2, that model changed:

- `textinput.Model` no longer exposes a public cursor state field. It stores a private `virtualCursor` at `charm.land/bubbles/v2@v2.1.0/textinput/textinput.go:104`.
- `textinput.Model.Cursor()` is now a method for returning a real terminal cursor, and it returns `nil` when the default virtual cursor is enabled at `charm.land/bubbles/v2@v2.1.0/textinput/textinput.go:916`.
- The `Blink` field exposed through `Cursor()` is only the static style/config flag copied from `styles.Cursor.Blink` at `charm.land/bubbles/v2@v2.1.0/textinput/textinput.go:932`. It is not the live blink phase.
- The actual live blink phase is private state, `IsBlinked`, inside the cursor model at `charm.land/bubbles/v2@v2.1.0/cursor/cursor.go:78`, and it toggles at `charm.land/bubbles/v2@v2.1.0/cursor/cursor.go:149`.

## What should be tested instead

Because superfile uses the default virtual cursor path, the meaningful public signal is rendered output:

- `textinput.Blink()` only seeds the blink cycle.
- The returned command produces the actual blink message.
- That message flips the internal blink state.
- The visible effect is in `cursor.Model.View()` at `charm.land/bubbles/v2@v2.1.0/cursor/cursor.go:235`.

So the correct v2 test is to:

1. Trigger the blink cycle.
2. Feed the resulting blink messages back into `HandleUpdate`.
3. Assert that `m.textInput.View()` changes and then changes back.

## Conclusion

Checking `m.textInput.Cursor().Blink` in v2 would not verify runtime blinking behavior:

- with virtual cursor enabled, `Cursor()` is `nil`
- with real cursor enabled, `.Blink` only reflects configuration, not the current blink phase

That is why the restored test uses rendered output instead of `Cursor().Blink`.
