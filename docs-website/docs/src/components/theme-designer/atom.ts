class Atom<T> {
	#v: T;
	#subscribers = new Map<(v: T) => void, (v: T) => void>();
	#notify = () => this.#subscribers.forEach((cb) => cb(this.#v));
	constructor(init: T) {
		this.#v = init;
	}
	get(): T {
		return this.#v;
	}
	set(v: T): void {
		this.#v = v;
		this.#notify();
	}
	subscribe(cb: (v: T) => void): () => boolean {
		cb(this.#v);
		this.#subscribers.set(cb, cb);
		return () => this.#subscribers.delete(cb);
	}
}

type MapStore<T> = Atom<T> & { setKey: (key: keyof T, value: T[typeof key]) => void };

export function map<T extends Record<string, unknown>>(value: T): MapStore<T> {
	const atom = new Atom(value) as MapStore<T>;
	atom.setKey = (key: keyof T, value: T[typeof key]) => {
		const curr = atom.get();
		if (curr[key] !== value) atom.set({ ...curr, [key]: value });
	};
	return atom;
}
