// The dot-matrix board behind the footer, drawn on a canvas that fills it.
// Every cell is a dot. Lettering dots are paper-coloured and sized by how much of the cell the
// letter covers; the rest are faint. "Energy" (0 to 1) makes a dot larger and orange: now and then
// a random dot of the lettering lights up and fades, and the pointer leaves a patch that fades too.

// "superfile" in IBM Plex Mono Bold, sampled onto a 96-column grid. Each digit is how much of that
// cell the lettering covers: 0 none, 1 a little, 2 most, 3 all.
const LETTERING = [
	"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	"000000000000000000000000000000000000000000000000000000000023333100002330000233332000000000000000",
	"000000000000000000000000000000000000000000000000000000000233333100003331000333333000000000000000",
	"000000000000000000000000000000000000000000000000000000000333111000003330000111333000000000000000",
	"000000000000000000000000000000000000000000000000000000000333000000000000000000333000000000000000",
	"000123332000222000222012211332000002333200002222201220222333222002222220000000333000000123321000",
	"002333333300333000333013323333300033333330013333323330333333333103333330000000333000002333333100",
	"003330013100333000333013331133320233101332002233333320122333222001223330000000333000003320023300",
	"003331000000333000333013310003320332000233000033310000000333000000003330000000333000023300013310",
	"003333332100333000333013310003330333333333000033300000000333000000003330000000333000023333333320",
	"000233333300333000333013310003330333222222000033300000000333000000003330000000333000023322222210",
	"000000123310333000333013310003320333000000000033300000000333000000003330000000333000023300000000",
	"003310023310333101333013331023320233201231001133311100111333111001113331110111333111013331013200",
	"013333333300133333333013323333300033333332013333333300233333333003333333320333333333002333333300",
	"001233332000023330232013312333100002333310003333333200233333332002333333320233333332000133332000",
	"000000000000000000000013310000000000000000000000000000000000000000000000000000000000000000000000",
	"000000000000000000000013310000000000000000000000000000000000000000000000000000000000000000000000",
	"000000000000000000000013310000000000000000000000000000000000000000000000000000000000000000000000",
	"000000000000000000000012200000000000000000000000000000000000000000000000000000000000000000000000",
	"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
];
const LETTER_COLS = LETTERING[0].length;

const INK = [239, 232, 221]; // term-ink
const DIM = [78, 88, 99]; // term-line
const HOT = [255, 87, 34]; // signal
const RADIUS = [0.09, 0.2, 0.31, 0.4]; // dot radius per coverage level, in cells
const SPARK_EVERY = 110; // a new lettering dot lights up this often, in milliseconds
const FADE_TAKES = 1400; // and takes this long to go back to paper
const POINTER_REACH = 4.5; // how far the pointer's glow spreads, in cells

const mix = (from: number[], to: number[], t: number) => `rgb(${from.map((value, i) => Math.round(value + (to[i] - value) * t)).join(",")})`;

export function mountBoard(canvas: HTMLCanvasElement) {
	const footer = canvas.parentElement!;
	const context = canvas.getContext("2d")!;
	const calm = matchMedia("(prefers-reduced-motion: reduce)");

	// The lettering is exactly as wide as the footer's content column and sits at the bottom, in the
	// space the footer leaves free under its links. The faint dots carry on to both edges of the page.
	let cell = 1; // one cell, in CSS pixels
	let originX = 0; // where column 0 starts; it may be off the left edge
	let colCount = 0;
	let rowCount = 0;
	let coverage: number[] = []; // per cell, 0 to 3
	let lettering: number[] = []; // the cells that belong to the lettering
	let glow = new Float32Array(0); // per cell, the energy left

	let visible = false;
	let frame = 0;
	let last = 0;
	let nextSpark = 0;

	function layout() {
		const column = footer.querySelector<HTMLElement>(".sf-wrap")!;
		const style = getComputedStyle(column);
		const left = column.offsetLeft + parseFloat(style.paddingLeft);
		const width = column.clientWidth - parseFloat(style.paddingLeft) - parseFloat(style.paddingRight);
		cell = width / LETTER_COLS;
		const colsBefore = Math.ceil(left / cell);
		originX = left - colsBefore * cell;
		colCount = colsBefore + LETTER_COLS + Math.ceil((canvas.clientWidth - left - width) / cell);
		rowCount = Math.max(LETTERING.length, Math.ceil(canvas.clientHeight / cell));
		const rowsAbove = rowCount - LETTERING.length - 1;
		coverage = [];
		for (let row = 0; row < rowCount; row++) {
			const line = LETTERING[row - rowsAbove];
			for (let col = 0; col < colCount; col++) coverage.push(Number(line?.[col - colsBefore] ?? 0));
		}
		lettering = coverage.flatMap((value, index) => (value >= 2 ? [index] : []));
		glow = new Float32Array(coverage.length);
	}

	function draw() {
		const ratio = canvas.width / canvas.clientWidth;
		const size = cell * ratio;
		context.clearRect(0, 0, canvas.width, canvas.height);
		for (let row = 0; row < rowCount; row++) {
			for (let col = 0; col < colCount; col++) {
				const index = row * colCount + col;
				const lit = coverage[index];
				const energy = glow[index];
				context.fillStyle = mix(lit ? INK : DIM, HOT, energy);
				context.globalAlpha = lit ? 1 : 0.4 + energy * 0.6;
				context.beginPath();
				context.arc(originX * ratio + (col + 0.5) * size, (row + 0.5) * size, (RADIUS[lit] + energy * (lit ? 0.06 : 0.14)) * size, 0, Math.PI * 2);
				context.fill();
			}
		}
		context.globalAlpha = 1;
	}

	function resize() {
		const ratio = Math.min(window.devicePixelRatio, 2);
		canvas.width = Math.round(canvas.clientWidth * ratio);
		canvas.height = Math.round(canvas.clientHeight * ratio);
		layout();
		draw();
	}

	function tick(now: number) {
		frame = 0;
		const elapsed = Math.min(64, now - last);
		last = now;
		let glowing = false;
		for (let i = 0; i < glow.length; i++) {
			if (glow[i] > 0) {
				glow[i] = Math.max(0, glow[i] - elapsed / FADE_TAKES);
				glowing ||= glow[i] > 0;
			}
		}
		if (!calm.matches && now >= nextSpark && lettering.length) {
			glow[lettering[Math.floor(Math.random() * lettering.length)]] = 1;
			nextSpark = now + SPARK_EVERY * (0.5 + Math.random());
			glowing = true;
		}
		draw();
		// With motion allowed there is always another spark coming; otherwise stop once the glow is gone.
		if (glowing || !calm.matches) run();
	}

	function run() {
		if (frame || !visible) return;
		last = performance.now();
		frame = requestAnimationFrame(tick);
	}

	// The canvas is behind the footer's content, so the footer itself reports the pointer.
	footer.addEventListener("pointermove", event => {
		const box = canvas.getBoundingClientRect();
		const x = (event.clientX - box.left - originX) / cell;
		const y = (event.clientY - box.top) / cell;
		for (let row = Math.max(0, Math.floor(y - POINTER_REACH)); row <= Math.min(rowCount - 1, Math.ceil(y + POINTER_REACH)); row++) {
			for (let col = Math.max(0, Math.floor(x - POINTER_REACH)); col <= Math.min(colCount - 1, Math.ceil(x + POINTER_REACH)); col++) {
				const near = 1 - Math.hypot(col + 0.5 - x, row + 0.5 - y) / POINTER_REACH;
				const index = row * colCount + col;
				if (near > glow[index]) glow[index] = near;
			}
		}
		run();
	});

	new ResizeObserver(resize).observe(canvas);
	new IntersectionObserver(([entry]) => {
		visible = entry.isIntersecting;
		run();
	}).observe(canvas);
}
