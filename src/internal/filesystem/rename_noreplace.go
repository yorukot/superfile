package filesystem

import "errors"

var errNoReplaceUnsupported = errors.New("atomic no-replace rename is unsupported")
