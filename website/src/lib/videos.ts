export interface Video {
	/** The YouTube video id. */
	id: string;
	title: string;
	creator: string;
	duration: string;
}

export const videos: Video[] = [
	{ id: "pDIMXiKyvA4", title: "Superfile is AWESOME", creator: "TechHut", duration: "12:38" },
	{ id: "3LEctgQWG10", title: "Superfile - AMAZING Terminal File Manager Created by a High School Kid", creator: "Better Stack", duration: "4:45" },
	{ id: "tYRlUQIPboU", title: "This Terminal File Manager is Actually Insane", creator: "Seth Phaeno", duration: "8:10" },
	{ id: "bCibCIpq5ko", title: "I wish I knew about this package sooner", creator: "Melkey", duration: "12:56" },
	{ id: "IGnbTs5HWvY", title: "Warp + Neovim + Superfile: This Terminal Setup CHANGED How I Code", creator: "Better Stack", duration: "4:31" },
	{ id: "LGipyvaecIU", title: "Este gestor de archivos es una auténtica LOCURA", creator: "atareao", duration: "24:42" },
	{ id: "cq4ndzcDsI4", title: "Superfile: Gerencie Arquivos pelo Terminal Facilmente!", creator: "Domynix Tech", duration: "13:59" },
	{ id: "DTbnVSapFEA", title: "Un FILE MANAGER sul TERMINALE...? - Superfile", creator: "YoungTech", duration: "7:00" },
	{ id: "7xtxeLjTS4Q", title: "Superfile: The Modern, Lightweight Linux Terminal File Manager", creator: "CodingBite", duration: "2:37" }
];

export const videoUrl = (id: string) => `https://www.youtube.com/watch?v=${id}`;
export const videoThumbnail = (id: string) => `https://i.ytimg.com/vi/${id}/maxresdefault.jpg`;
