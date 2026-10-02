//go:build darwin && cgo

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
	char *data;
	size_t length;
	char *errorMessage;
} SPFClipResult;

// spf_path_from_cstr decodes a filesystem path using macOS conventions, returning
// nil for a NULL pointer or a path that cannot be represented as an NSString.
static NSString *spf_path_from_cstr(const char *path) {
	if (path == NULL) {
		return nil;
	}
	return [[NSFileManager defaultManager]
		stringWithFileSystemRepresentation:path
		length:strlen(path)];
}

// spf_clipboard_copy_files writes the given file paths to the general pasteboard
// as file URLs. Returns an error string (caller frees) or NULL on success.
char *spf_clipboard_copy_files(const char **paths, int count) {
	@autoreleasepool {
		if (paths == NULL || count <= 0) {
			return strdup("no files to copy");
		}
		NSMutableArray<NSURL *> *urls = [NSMutableArray arrayWithCapacity:count];
		for (int i = 0; i < count; i++) {
			NSString *p = spf_path_from_cstr(paths[i]);
			if (p == nil) {
				return strdup("failed to create macOS file path");
			}
			[urls addObject:[NSURL fileURLWithPath:p]];
		}
		NSPasteboard *pb = [NSPasteboard generalPasteboard];
		[pb clearContents];
		if (![pb writeObjects:urls]) {
			return strdup("failed to write file URLs to pasteboard");
		}
		return NULL;
	}
}

// spf_clipboard_paste_files returns the file paths currently on the general
// pasteboard in their filesystem representation, each terminated by a NUL byte.
// NUL cannot occur in a macOS path, so this framing is lossless even for names
// containing newlines or trailing whitespace. data is NULL and length is 0 when
// the pasteboard holds no file URLs. Caller frees data and errorMessage.
SPFClipResult spf_clipboard_paste_files(void) {
	SPFClipResult result = {0};
	@autoreleasepool {
		NSPasteboard *pb = [NSPasteboard generalPasteboard];
		NSDictionary *options = @{ NSPasteboardURLReadingFileURLsOnlyKey : @YES };
		NSArray *classes = @[ [NSURL class] ];
		NSArray<NSURL *> *urls = [pb readObjectsForClasses:classes options:options];

		NSMutableData *buf = [NSMutableData data];
		for (NSURL *url in urls) {
			if (![url isFileURL]) {
				continue;
			}
			const char *fsPath = [url fileSystemRepresentation];
			if (fsPath == NULL) {
				continue;
			}
			// Include the terminating NUL as the record separator.
			[buf appendBytes:fsPath length:strlen(fsPath) + 1];
		}

		if ([buf length] == 0) {
			return result;
		}
		result.data = malloc([buf length]);
		if (result.data == NULL) {
			result.errorMessage = strdup("out of memory reading pasteboard");
			return result;
		}
		memcpy(result.data, [buf bytes], [buf length]);
		result.length = [buf length];
		return result;
	}
}

// spf_clip_free_string releases a string returned by the clipboard bridge;
// passing NULL is safe.
void spf_clip_free_string(char *value) {
	if (value != NULL) {
		free(value);
	}
}
