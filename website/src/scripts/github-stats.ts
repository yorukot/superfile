// Numbers about the repository, asked from GitHub in the browser and kept for an hour.
import { REPO } from "../lib/site";

const API = `https://api.github.com/repos/${REPO}`;
const MAX_AGE = 3_600_000;

/** Runs `load` unless a fresh answer is stored. Null when GitHub cannot be reached or storage is blocked. */
async function cached(key: string, load: () => Promise<number>): Promise<number | null> {
	try {
		const stored = JSON.parse(localStorage.getItem(key) ?? "null") as { n: number; at: number } | null;
		if (stored && Date.now() - stored.at < MAX_AGE) return stored.n;
		const n = await load();
		localStorage.setItem(key, JSON.stringify({ n, at: Date.now() }));
		return n;
	} catch {
		return null;
	}
}

async function get(path = ""): Promise<Response> {
	const response = await fetch(API + path);
	if (!response.ok) throw new Error(response.statusText);
	return response;
}

export const getStars = () => cached("spf-stars", async () => (await (await get()).json()).stargazers_count);

// Asking for one contributor per page makes the number of the last page the total.
export const getContributors = () =>
	cached("spf-contributors", async () => {
		const link = (await get("/contributors?per_page=1&anon=true")).headers.get("link") ?? "";
		return Number(/[?&]page=(\d+)>; rel="last"/.exec(link)?.[1] ?? 1);
	});

// Every asset's download count across all releases. Installs through package managers are not counted.
export const getDownloads = () =>
	cached("spf-downloads", async () => {
		const releases: { assets: { download_count: number }[] }[] = await (await get("/releases?per_page=100")).json();
		return releases.flatMap(release => release.assets).reduce((sum, asset) => sum + asset.download_count, 0);
	});
