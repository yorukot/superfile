# PR 1113 Outstanding Review Items

## 1. Bulk rename text input width alignment
- **Commenter**: CodeRabbit (auto-review)
- **File**: `src/internal/common/style_function.go`
- **Issue**: `GenerateBulkRenameTextInput` hard-codes `ti.Width = ModalWidth - 10`, so the bulk rename inputs overflow the 54-column modal, hiding the text cursor.
- **Required Action**:
  - Update `GenerateBulkRenameTextInput` to accept a `width int` parameter and use it for `ti.Width`.
  - Pass the modal column width (54) from every call site in `src/internal/ui/bulk_rename/model.go`.
  - Confirm the rendered input fields stay within the right column and the cursor remains visible during typing.

## 2. No other actionable reviewer feedback
- Remaining review threads are informational or duplicates; only item 1 is pending.
