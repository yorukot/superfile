//go:build darwin && cgo

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
	char *data;
	char *errorMessage;
} SPFClipResult;

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
// pasteboard, newline-separated. The result is empty (but non-NULL) when the
// pasteboard holds no file URLs. Caller frees both fields.
SPFClipResult spf_clipboard_paste_files(void) {
	SPFClipResult result = {0};
	@autoreleasepool {
		NSPasteboard *pb = [NSPasteboard generalPasteboard];
		NSDictionary *options = @{ NSPasteboardURLReadingFileURLsOnlyKey : @YES };
		NSArray *classes = @[ [NSURL class] ];
		NSArray<NSURL *> *urls = [pb readObjectsForClasses:classes options:options];

		NSMutableString *joined = [NSMutableString string];
		for (NSURL *url in urls) {
			if (![url isFileURL]) {
				continue;
			}
			if ([joined length] > 0) {
				[joined appendString:@"\n"];
			}
			[joined appendString:[url path]];
		}
		result.data = strdup([joined UTF8String]);
		return result;
	}
}

void spf_clip_free_string(char *value) {
	if (value != NULL) {
		free(value);
	}
}
