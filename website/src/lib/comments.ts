// The GitHub discussion that holds the comments of each docs page, in the "Docs Comments" category.
//
// This table is a workaround. giscus normally finds a page's discussion by searching GitHub for its
// title, but GitHub's discussion search index for this repository is incomplete (it knows 3 of 81
// discussions), so that search reports "0 comments" for threads that exist. The same symptom is
// described in https://github.com/giscus/giscus/issues/1268. Pointing at the discussion by number
// does not go through search.
//
// A page that is not listed here still works the old way: giscus searches by title, and creates the
// discussion on the first comment. Once GitHub's index is repaired this file can be deleted.
export const discussionNumbers: Record<string, number> = {
	overview: 167,
	changelog: 1711,
	troubleshooting: 549,
	"getting-started/installation": 175,
	"getting-started/tutorial": 185,
	"getting-started/image-preview": 939,
	"configure/config-file-path": 184,
	"configure/superfile-config": 959,
	"configure/custom-hotkeys": 285,
	"configure/custom-theme": 1071,
	"list/hotkey-list": 319,
	"list/theme-list": 186,
	"list/plugin-list": 1045
};
