// Releases are read from the changelog page at build time, so the home page never repeats something that can go stale.
import fs from "node:fs";

export interface Release {
	version: string;
	date: string;
	url: string;
	/** The entries under "Update", without their pull request links. */
	highlights: string[];
}

function parseChangelog(source: string): Release[] {
	const releases: Release[] = [];
	const heading = /^# \[\*\*([^*]+)\*\*\]\(([^)]+)\)/;
	let release: Release | undefined;
	let section = "";

	for (const line of source.split("\n")) {
		const match = heading.exec(line);
		if (match) {
			release = { version: match[1], url: match[2], date: "", highlights: [] };
			releases.push(release);
			section = "";
		} else if (!release) {
			continue;
		} else if (line.startsWith("> ") && !release.date) {
			release.date = line.slice(2).trim();
		} else if (line.startsWith("#### ")) {
			section = line.slice(5).trim().toLowerCase();
		} else if (line.startsWith("- ") && section === "update") {
			const text = line
				.slice(2)
				.replace(/\s*\[`#\d+`\]\([^)]*\)/g, "")
				.replace(/\s+by\s+@\S+$/i, "")
				.trim();
			if (text) release.highlights.push(text);
		}
	}
	return releases;
}

/** Every release in the changelog, newest first. */
export const releases = parseChangelog(fs.readFileSync("src/content/docs/changelog.md", "utf-8"));
