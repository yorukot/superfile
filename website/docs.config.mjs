// How the docs are versioned. Both the build script (scripts/sync-versions.mjs) and the pages read this.
//
// The version switcher offers:
//   latest   the docs in this branch (src/content/docs). They follow development and have no version
//            number, so they can describe things that are not released yet.
//   v1.6.0…  the docs of every release and release candidate, read from its git tag.
export const docsConfig = {
	// Which docs open at the site root (/overview, /getting-started/…):
	//   "newest"  the newest release (the default: readers see the docs of what they can install)
	//   "latest"  the docs in this branch
	//   a tag     one particular release, for example "v1.6.0"
	// Everything else is reachable under its own name: /latest/…, /v1.5.0/…
	root: "newest"
};

/** The name the docs in this branch go by, in URLs and in the version switcher. */
export const LATEST = "latest";

/** Reads a tag such as "v1.4.0", "v1.4.1-rc", "v1.6.0-rc1" or "v1.2.0.0". Returns null for anything else. */
export function parseTag(tag) {
	const match = /^v(\d+)\.(\d+)\.(\d+)(?:\.(\d+))?(?:-(rc\d*))?$/.exec(tag);
	if (!match) return null;
	const [major, minor, patch, build] = match.slice(1, 5).map(part => Number(part ?? 0));
	return { major, minor, patch, build, pre: match[5] ?? "", minorKey: `v${major}.${minor}` };
}

const isRelease = tag => !parseTag(tag)?.pre;

/** Orders tags newest first. A release candidate comes after the release it led to. */
export function compareTags(a, b) {
	const [x, y] = [parseTag(a), parseTag(b)];
	if (!x || !y) return a < b ? 1 : -1;
	return y.major - x.major || y.minor - x.minor || y.patch - x.patch || y.build - x.build || (x.pre ? 1 : 0) - (y.pre ? 1 : 0) || (x.pre < y.pre ? 1 : -1);
}

/**
 * Works out, from the release tags that could be read, where every set of docs lives.
 * Without any tag (a build with no git history) the docs in this branch are served at the root.
 * @param {string[]} found
 */
export function resolveDocs(found) {
	const tags = [...found].sort(compareTags);
	// The newest real release: release candidates do not count.
	const newest = tags.find(isRelease) ?? tags[0];
	const wanted = docsConfig.root === "newest" ? newest : docsConfig.root;
	const root = wanted && tags.includes(wanted) ? wanted : LATEST;

	// The switcher shows one line per minor version (its newest release); the patch releases and
	// release candidates before it are folded underneath.
	const byMinor = Map.groupBy(tags, tag => parseTag(tag)?.minorKey ?? tag);
	const groups = Array.from(byMinor, ([key, members]) => {
		const head = members.find(isRelease) ?? members[0];
		return { key, head, rest: members.filter(tag => tag !== head) };
	});

	return {
		/** The docs at the site root: a tag, or "latest". */
		root,
		/** The newest release, if any. */
		newest,
		/** The releases by minor version, newest first: { key: "v1.3", head: "v1.3.3", rest: ["v1.3.2", …] }. */
		groups,
		/** The URL prefix of a set of docs: "" at the root, "/latest" or "/v1.5.0" elsewhere. */
		prefix: version => (version === root ? "" : `/${version}`)
	};
}
