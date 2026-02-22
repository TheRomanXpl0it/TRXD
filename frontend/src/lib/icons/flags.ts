let loaded = false;

export async function ensureCircleFlags() {
	if (loaded) return;
	const [{ addCollection }, data] = await Promise.all([
		import('@iconify/svelte'),
		import('@iconify-json/circle-flags/icons.json')
	]);
	// data.default for ESM json, fallback otherwise
	addCollection((data as any).default ?? data);
	loaded = true;
}
