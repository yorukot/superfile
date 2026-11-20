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
const DENSITY = 1.5; // canvas pixels per CSS pixel, at most: the board costs its area again on every update
const SPARK_EVERY = 110; // a new lettering dot lights up this often, in milliseconds
const FADE_TAKES = 1400; // and takes this long to go back to paper
const STEP = 60; // the fade is slow: it moves this often, in milliseconds, not on every frame
const POINTER_STEP = 30; // and this often while the pointer is drawing on the board
const POINTER_REACH = 4.5; // how far the pointer's glow spreads, in cells

const mix = (from: number[], to: number[], t: number) => `rgb(${from.map((value, i) => Math.round(value + (to[i] - value) * t)).join(",")})`;

export function mountBoard(canvas: HTMLCanvasElement) {
	const footer = canvas.parentElement!;
	const context = canvas.getContext("2d")!;
	const calm = matchMedia("(prefers-reduced-motion: reduce)");

	// The lettering is exactly as wide as the footer's content column and sits at the bottom, in the
	// space the footer leaves free under its links. The faint dots carry on to both edges of the page.
	let ratio = 1; // canvas pixels per CSS pixel
	let cell = 1; // one cell, in CSS pixels
	let originX = 0; // where column 0 starts; it may be off the left edge
	let colCount = 0;
	let rowCount = 0;
	let coverage: number[] = []; // per cell, 0 to 3
	let lettering: number[] = []; // the cells that belong to the lettering
	let glow = new Float32Array(0); // per cell, the energy left
	const lit = new Set<number>(); // the cells with energy left: the only ones a step repaints

	let visible = false; // the footer is on screen
	let moved = false; // the pointer lit cells since the last step
	let timer = 0;
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
		lit.clear();
	}

	// Nothing runs while the footer is off screen or the tab is in the background.
	const showing = () => visible && !document.hidden;

	const radius = (index: number) => RADIUS[coverage[index]] + glow[index] * (coverage[index] ? 0.06 : 0.14); // in cells

	function dot(index: number) {
		const size = cell * ratio;
		const col = index % colCount;
		const level = coverage[index];
		const energy = glow[index];
		context.fillStyle = mix(level ? INK : DIM, HOT, energy);
		context.globalAlpha = level ? 1 : 0.4 + energy * 0.6;
		context.beginPath();
		context.arc(originX * ratio + (col + 0.5) * size, ((index - col) / colCount + 0.5) * size, radius(index) * size, 0, Math.PI * 2);
		context.fill();
	}

	// Paints one cell again. The soft edge of a dot reaches up to a pixel past its cell, so the square
	// that is cleared is a pixel larger all round, and the dots that reach into it are drawn again too.
	function repaint(index: number) {
		const size = cell * ratio;
		const col = index % colCount;
		const row = (index - col) / colCount;
		const x = Math.floor(originX * ratio + col * size) - 1;
		const y = Math.floor(row * size) - 1;
		const width = Math.ceil(originX * ratio + (col + 1) * size) + 1 - x;
		const height = Math.ceil((row + 1) * size) + 1 - y;
		context.save();
		context.beginPath();
		context.rect(x, y, width, height);
		context.clip();
		context.clearRect(x, y, width, height);
		for (let r = Math.max(0, row - 1); r <= Math.min(rowCount - 1, row + 1); r++) {
			for (let c = Math.max(0, col - 1); c <= Math.min(colCount - 1, col + 1); c++) {
				const other = r * colCount + c;
				if (other === index || radius(other) > 0.5 - 3.5 / size) dot(other);
			}
		}
		context.restore();
	}

	function resize() {
		const density = Math.min(window.devicePixelRatio, DENSITY);
		canvas.width = Math.round(canvas.clientWidth * density);
		canvas.height = Math.round(canvas.clientHeight * density);
		ratio = canvas.width / canvas.clientWidth;
		layout();
		for (let index = 0; index < coverage.length; index++) dot(index);
	}

	function step() {
		const now = performance.now();
		const elapsed = Math.min(STEP, now - last);
		last = now;
		for (const index of lit) glow[index] = Math.max(0, glow[index] - elapsed / FADE_TAKES);
		if (!calm.matches && now >= nextSpark && lettering.length) {
			const index = lettering[Math.floor(Math.random() * lettering.length)];
			glow[index] = 1;
			lit.add(index);
			// Counted from when it was due, so that the rhythm does not depend on the step.
			nextSpark = Math.max(nextSpark, now - STEP) + SPARK_EVERY * (0.5 + Math.random());
		}
		for (const index of lit) {
			repaint(index);
			if (!glow[index]) lit.delete(index);
		}
		// With motion allowed there is always another spark coming; otherwise stop once the glow is gone.
		timer = showing() && (lit.size || !calm.matches) ? window.setTimeout(step, moved ? POINTER_STEP : STEP) : 0;
		moved = false;
	}

	function run() {
		if (timer || !showing()) return;
		timer = window.setTimeout(step);
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
				if (near > glow[index]) {
					glow[index] = near;
					lit.add(index);
					moved = true;
				}
			}
		}
		run();
	});

	document.addEventListener("visibilitychange", run);
	new ResizeObserver(resize).observe(canvas);
	new IntersectionObserver(([entry]) => {
		visible = entry.isIntersecting;
		run();
	}).observe(canvas);
}
