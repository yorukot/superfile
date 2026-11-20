import type { BrandIconName } from "./brand-icons";

export type OS = "linux" | "macos" | "windows";

export const osNames: Record<OS, string> = { linux: "Linux", macos: "macOS", windows: "Windows" };

export interface InstallMethod {
	label: string;
	icon: BrandIconName;
	prompt: string;
	cmd: string;
	/** The systems this method is for. Within one system, the first method listed is the one the page leads with. */
	os: OS[];
}

/** What the page shows before it knows the visitor's system. */
export const defaultInstallMethod: InstallMethod = {
	label: "linux / macos",
	icon: "linux",
	prompt: "$",
	cmd: 'bash -c "$(curl -sLo- https://superfile.dev/install.sh)"',
	os: ["linux", "macos"]
};

export const installMethods: InstallMethod[] = [
	{ label: "homebrew", icon: "homebrew", prompt: "$", cmd: "brew install superfile", os: ["macos"] },
	defaultInstallMethod,
	{ label: "winget", icon: "windows", prompt: "PS>", cmd: "winget install --id yorukot.superfile", os: ["windows"] },
	{ label: "arch", icon: "archlinux", prompt: "$", cmd: "sudo pacman -S superfile", os: [] },
	{ label: "nix", icon: "nixos", prompt: "$", cmd: "nix profile install github:yorukot/superfile#superfile", os: [] }
];

/** The visitor's system, in the browser. Phones and tablets get none: nothing here installs on them. */
export function detectOS(): OS | null {
	const platform = `${navigator.platform} ${navigator.userAgent}`.toLowerCase();
	if (/android|iphone|ipad/.test(platform)) return null;
	if (platform.includes("win")) return "windows";
	if (platform.includes("mac")) return "macos";
	return /linux|x11|bsd/.test(platform) ? "linux" : null;
}
