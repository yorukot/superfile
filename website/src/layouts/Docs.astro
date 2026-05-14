---
import type { CollectionEntry } from "astro:content";
import { getEntry } from "astro:content";
import type { MarkdownHeading } from "astro";
import Giscus from "../components/Giscus.astro";
import SharedHead from "../components/SharedHead.astro";
import SiteNav from "../components/SiteNav.astro";
import "../styles/global.css";
import "../styles/docs.css";
import { navigation, isGroup, getPrevNext, type NavLeaf } from "../lib/navigation";

interface Props {
	entry: CollectionEntry<"docs">;
	headings: MarkdownHeading[];
}

const { entry, headings } = Astro.props;
const h1Heading = headings.find(h => h.depth === 1);
const title = entry.data.title ?? h1Heading?.text ?? "superfile";
const { description } = entry.data;

// Determine locale and base slug
// Astro v5 glob loader: entry.id has no file extension
const isZhTw = entry.id.startsWith("zh-tw/");
const baseSlug = isZhTw ? entry.id.replace(/^zh-tw\//, "") : entry.id;
const locale = isZhTw ? "zh-tw" : "";
const localePrefix = isZhTw ? "/zh-tw" : "";
const enEntry = isZhTw ? await getEntry("docs", baseSlug) : entry;
const giscusTerm = `${enEntry?.data.title ?? title} | superfile`;

// Current path for active detection
const currentPath = Astro.url.pathname.replace(/\/$/, "");

function isActive(slug: string) {
	const path = `${localePrefix}/${slug}`;
	return currentPath === path || currentPath === path + "/";
}

// Prev / next
const { prev, next } = getPrevNext(baseSlug);

// TOC: only h2 and h3
const tocHeadings = headings.filter(h => h.depth === 2 || h.depth === 3);

// If the content body already opens with an h1, skip rendering the layout h1
const hasContentH1 = headings.length > 0 && headings[0].depth === 1;

// Language toggle target
function langTogglePath(slug: string) {
	if (isZhTw) return `/${slug}`;
	return `/zh-tw/${slug}`;
}

const siteUrl = "https://superfile.dev";
const canonicalUrl = new URL(Astro.url.pathname, siteUrl).href;
---

<!doctype html>
<html lang={isZhTw ? "zh-TW" : "en"}>
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<title>{title} | superfile</title>
		{description && <meta name="description" content={description} />}
		<link rel="canonical" href={canonicalUrl} />
		<meta property="og:title" content={`${title} | superfile`} />
		{description && <meta property="og:description" content={description} />}
		<meta property="og:image" content={`${siteUrl}/og.png?v=1`} />
		<meta name="twitter:card" content="summary_large_image" />

		<SharedHead />
	</head>

	<body>
		<SiteNav isDocs={true} isZhTw={isZhTw} baseSlug={baseSlug} />

		<!-- ── Overlay (mobile) ──────────────────────────────────────────────── -->
		<div class="docs-overlay" id="docs-overlay"></div>

		<!-- ── Shell ─────────────────────────────────────────────────────────── -->
		<div class="docs-shell">
			<!-- Sidebar -->
			<aside class="docs-sidebar" id="docs-sidebar">
				<nav class="docs-nav" aria-label="Docs navigation">
					{
						navigation.map(node => {
							if (isGroup(node)) {
								const label = isZhTw && node.zhLabel ? node.zhLabel : node.label;
								const hasActive = node.items.some(item => isActive(item.slug));
								return (
									<details class="docs-nav-section" open={hasActive || undefined}>
										<summary>
											<svg class="docs-nav-section-arrow" viewBox="0 0 12 12" fill="none">
												<path d="M4 2l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
											</svg>
											{label}
										</summary>
										<ul class="docs-nav-items">
											{node.items.map(item => {
												const itemLabel = isZhTw && item.zhLabel ? item.zhLabel : item.label;
												const href = `${localePrefix}/${item.slug}`;
												return (
													<li class={`docs-nav-item${isActive(item.slug) ? " docs-nav-item--active" : ""}`}>
														<a href={href}>{itemLabel}</a>
													</li>
												);
											})}
										</ul>
									</details>
								);
							} else {
								const label = isZhTw && node.zhLabel ? node.zhLabel : node.label;
								const href = `${localePrefix}/${node.slug}`;
								return (
									<a href={href} class={`docs-nav-leaf${isActive(node.slug) ? " docs-nav-leaf--active" : ""}`}>
										{label}
									</a>
								);
							}
						})
					}
				</nav>
			</aside>

			<!-- Body: centered content + TOC -->
			<div class="docs-body">
				<div class="docs-body-inner">
					<!-- Main -->
					<main class="docs-main">
						<article class="docs-content" data-pagefind-body data-pagefind-meta={`title:${title}`}>
							{
								!hasContentH1 && (
									<div class="docs-page-header">
										<h1>{title}</h1>
										{description && <p class="docs-description">{description}</p>}
									</div>
								)
							}
							<slot />
						</article>

						<!-- Prev / Next -->
						{
							(prev || next) && (
								<nav class="docs-page-nav" aria-label="Page navigation">
									{prev ? (
										<a href={`${localePrefix}/${prev.slug}`} class="docs-page-nav-link docs-page-nav-link--prev">
											<span class="docs-page-nav-label">← Previous</span>
											<span class="docs-page-nav-title">{isZhTw && prev.zhLabel ? prev.zhLabel : prev.label}</span>
										</a>
									) : (
										<div />
									)}
									{next ? (
										<a href={`${localePrefix}/${next.slug}`} class="docs-page-nav-link docs-page-nav-link--next">
											<span class="docs-page-nav-label">Next →</span>
											<span class="docs-page-nav-title">{isZhTw && next.zhLabel ? next.zhLabel : next.label}</span>
										</a>
									) : (
										<div />
									)}
								</nav>
							)
						}

						<!-- Giscus comments -->
						<div class="docs-giscus">
							<Giscus term={giscusTerm} lang={isZhTw ? "zh-TW" : "en"} />
						</div>
					</main>

					<!-- TOC -->
					{
						tocHeadings.length > 0 && (
							<aside class="docs-toc" aria-label="Table of contents">
								<p class="docs-toc-title">On this page</p>
								<ul class="docs-toc-list" id="docs-toc-list">
									{tocHeadings.map(h => (
										<li class={`docs-toc-item docs-toc-item--h${h.depth}`} data-heading={h.slug}>
											<a href={`#${h.slug}`}>{h.text}</a>
										</li>
									))}
								</ul>
							</aside>
						)
					}
				</div><!-- /.docs-body-inner -->
			</div><!-- /.docs-body -->
		</div><!-- /.docs-shell -->

		<script>
			// ── Module-scoped state (persists across SPA navigations) ─────────────
			let _tocObserver: IntersectionObserver | null = null;

			// Restore body scroll on every navigation (sidebar may have locked it)
			document.addEventListener("astro:before-swap", () => {
				document.body.style.overflow = "";
			});

			// Re-initialise all docs interactivity after every page load
			document.addEventListener("astro:page-load", () => {
				// ── Mobile sidebar toggle ───────────────────────────────────────────
				const menuBtn = document.getElementById("docs-menu-btn");
				const sidebar = document.getElementById("docs-sidebar");
				const overlay = document.getElementById("docs-overlay");

				function openSidebar() {
					sidebar?.classList.add("docs-sidebar--open");
					overlay?.classList.add("docs-overlay--visible");
					document.body.style.overflow = "hidden";
				}

				function closeSidebar() {
					sidebar?.classList.remove("docs-sidebar--open");
					overlay?.classList.remove("docs-overlay--visible");
					document.body.style.overflow = "";
				}

				menuBtn?.addEventListener("click", () => {
					sidebar?.classList.contains("docs-sidebar--open") ? closeSidebar() : openSidebar();
				});

				overlay?.addEventListener("click", closeSidebar);

				// ── Copy buttons ────────────────────────────────────────────────────
				document.querySelectorAll("pre").forEach(pre => {
					// Avoid double-adding on re-renders
					if (pre.querySelector(".docs-copy-btn")) return;
					const btn = document.createElement("button");
					btn.className = "docs-copy-btn";
					btn.textContent = "copy";
					btn.addEventListener("click", () => {
						const code = pre.querySelector("code");
						if (!code) return;
						navigator.clipboard.writeText(code.innerText).then(() => {
							btn.textContent = "copied!";
							setTimeout(() => (btn.textContent = "copy"), 2000);
						});
					});
					pre.appendChild(btn);
				});

				// ── TOC active state (IntersectionObserver) ─────────────────────────
				_tocObserver?.disconnect();
				_tocObserver = null;

				const tocItems = document.querySelectorAll<HTMLElement>("[data-heading]");
				if (tocItems.length) {
					const headingEls = Array.from(tocItems)
						.map(item => document.getElementById(item.dataset.heading!))
						.filter(Boolean) as HTMLElement[];

					_tocObserver = new IntersectionObserver(
						entries => {
							entries.forEach(entry => {
								const id = entry.target.id;
								const tocItem = document.querySelector(`[data-heading="${id}"]`);
								if (entry.isIntersecting) {
									document.querySelectorAll(".docs-toc-item--active").forEach(el => el.classList.remove("docs-toc-item--active"));
									tocItem?.classList.add("docs-toc-item--active");
								}
							});
						},
						{ rootMargin: "-20% 0% -60% 0%" }
					);

					headingEls.forEach(el => _tocObserver!.observe(el));
				}
			});
		</script>
	</body>
</html>
