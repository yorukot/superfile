// Copies the docs of every released version out of git, so the site can serve them next to the
// current docs: /v1.5.0/getting-started/tutorial and so on.
//
// For every release tag that contains the docs (v1.6.0, v1.4.1-rc, v1.3.2, …), this
// writes website/.versions/<tag>/ with the Markdown, the images it uses and the repository files
// it embeds, all read from that tag. Nothing here is committed: the folder is rebuilt from git.
//
// On your machine it never fails the build: without git, or without tags, it warns and the site is
// built with the current docs only. On a build service (CI is set) no tags at all is an error.
import { execFileSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { compareTags, parseTag } from "../docs.config.mjs";

const websiteDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const outDir = path.join(websiteDir, ".versions");
const DOCS = "website/src/content/docs";
// The component that embeds a repository file. Pages import it by a relative path, which differs in a copy.
const codeComponent = path.join(websiteDir, "src/components/code.astro");
// Raise this when the way a version is written out changes, so existing copies are rebuilt.
const FORMAT = 5;

// Paths given to git below are relative to the repository root, so git runs from there.
let gitDir = websiteDir;
const git = args => execFileSync("git", args, { cwd: gitDir, stdio: ["ignore", "pipe", "ignore"], maxBuffer: 256 * 1024 * 1024 });
const gitText = args => git(args).toString("utf-8");
const tryGit = args => {
	try {
		return gitText(args);
	} catch {
		return null;
	}
};

/** Every release and release-candidate tag, newest first. */
function releaseTags() {
	return (tryGit(["tag", "--list", "v*"]) ?? "")
		.split("\n")
		.filter(tag => parseTag(tag))
		.sort(compareTags);
}

/** Reads a file from a tag. A folder on the way may be a symlink there, which `git show` does not follow. */
function readFromTag(tag, file) {
	const parts = file.split("/");
	for (let depth = 1; depth <= parts.length; depth++) {
		const prefix = parts.slice(0, depth).join("/");
		const entry = tryGit(["ls-tree", tag, "--", prefix]);
		if (!entry) return null;
		if (entry.startsWith("120000")) {
			const target = gitText(["show", `${tag}:${prefix}`]).trim();
			const resolved = path.posix.normalize(path.posix.join(path.posix.dirname(prefix), target, ...parts.slice(depth)));
			return resolved.startsWith("..") ? null : readFromTag(tag, resolved);
		}
	}
	try {
		return git(["show", `${tag}:${file}`]);
	} catch {
		return null;
	}
}

/** Writes a file, creating the folders on the way. */
function write(file, content) {
	fs.mkdirSync(path.dirname(file), { recursive: true });
	fs.writeFileSync(file, content);
}

function syncVersion(tag) {
	const versionDir = path.join(outDir, tag);
	const stamp = `${gitText(["rev-parse", `${tag}^{commit}`]).trim()} format ${FORMAT}`;
	const stampFile = path.join(versionDir, ".stamp");
	if (fs.existsSync(stampFile) && fs.readFileSync(stampFile, "utf-8") === stamp) return "kept";

	const pages = gitText(["ls-tree", "-r", "--name-only", tag, "--", DOCS])
		.split("\n")
		.filter(file => /\.mdx?$/.test(file))
		.map(file => file.slice(DOCS.length + 1))
		// The translation is gone from the site, and index.mdx was the old home page.
		.filter(file => !file.startsWith("zh-tw/") && file !== "index.mdx");
	if (!pages.length) return "has no docs";

	fs.rmSync(versionDir, { recursive: true, force: true });
	let written = 0;

	for (const page of pages) {
		// A page can be a symlink in the tag (notice.md is one), so it is read the way that follows them.
		let text = readFromTag(tag, `${DOCS}/${page}`)?.toString("utf-8");
		// Pages built from components this site no longer has cannot be rendered.
		if (!text || /from\s+["']@astrojs\/starlight/.test(text)) continue;
		const target = path.join(versionDir, "src/content/docs", page);

		// Images: copy the ones the page uses to the same relative place; drop the ones the tag lacks.
		text = text.replace(/!\[([^\]]*)\]\((\.{1,2}\/[^)\s]+)\)/g, (whole, alt, relative) => {
			const inRepo = path.posix.normalize(path.posix.join(path.posix.dirname(`${DOCS}/${page}`), relative));
			const image = readFromTag(tag, inRepo);
			const copy = path.resolve(path.dirname(target), relative);
			// A path that climbs out of the copy's own folder would write over something else.
			if (!image || !copy.startsWith(versionDir + path.sep)) return "";
			write(copy, image);
			return whole;
		});

		// Embedded repository files (config.toml and friends) are read from this tag as well.
		if (/<CodeBlock\s/.test(text)) {
			const relative = path.relative(path.dirname(target), codeComponent).split(path.sep).join("/");
			text = text.replace(/^import CodeBlock from .*$/m, `import CodeBlock from "${relative}";`);
			text = text.replace(/<CodeBlock\s/g, `<CodeBlock root=".versions/${tag}/repo" `);
			for (const [, file] of text.matchAll(/file=["']([^"']+)["']/g)) {
				const embedded = readFromTag(tag, file);
				if (embedded) write(path.join(versionDir, "repo", file), embedded);
			}
		}

		write(target, text);
		written++;
	}

	fs.writeFileSync(stampFile, stamp);
	return `${written} pages`;
}

function main() {
	const root = tryGit(["rev-parse", "--show-toplevel"]);
	if (root === null) {
		console.warn("[versions] not a git checkout: building the current docs only");
		return;
	}
	gitDir = root.trim();
	// A shallow clone (most build services) lacks tags. One commit per tag is all that is needed.
	// This is only ever done to a clone that is already shallow.
	if (tryGit(["rev-parse", "--is-shallow-repository"])?.trim() === "true") tryGit(["fetch", "--depth=1", "--tags", "--quiet", "origin"]);
	const kept = [];
	for (const tag of releaseTags()) {
		if (!tryGit(["ls-tree", "-d", tag, "--", DOCS])) continue;
		try {
			console.log(`[versions] ${tag}: ${syncVersion(tag)}`);
			kept.push(tag);
		} catch (error) {
			console.warn(`[versions] ${tag}: skipped (${error.message.split("\n")[0]})`);
		}
	}
	// A copy whose tag no longer exists goes away.
	for (const existing of fs.existsSync(outDir) ? fs.readdirSync(outDir) : []) {
		if (!kept.includes(existing)) fs.rmSync(path.join(outDir, existing), { recursive: true, force: true });
	}
	if (kept.length) return;
	// A site deployed without its versions breaks every /v… link, so a build service stops here.
	if (process.env.CI && !process.env.DOCS_WITHOUT_VERSIONS) {
		console.error("[versions] no release tags found. Clone with the tags, or set DOCS_WITHOUT_VERSIONS=1 to build the current docs only.");
		process.exit(1);
	}
	console.warn("[versions] no release tags found: building the current docs only");
}

main();
