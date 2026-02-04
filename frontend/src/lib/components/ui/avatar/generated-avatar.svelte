<script lang="ts">
	let {
		seed,
		size = 120,
		class: className = ''
	} = $props<{ seed: string; size?: number; class?: string }>();

	// --- 1. Deterministic Randomness ---
	function cyrb128(input: any) {
		const str = String(input ?? '');
		let h1 = 1779033703,
			h2 = 3144134277,
			h3 = 1013904242,
			h4 = 2773480762;
		for (let i = 0, k; i < str.length; i++) {
			k = str.charCodeAt(i);
			h1 = h2 ^ Math.imul(h1 ^ k, 597399067);
			h2 = h3 ^ Math.imul(h2 ^ k, 2869860233);
			h3 = h4 ^ Math.imul(h3 ^ k, 951274213);
			h4 = h1 ^ Math.imul(h4 ^ k, 2716044179);
		}
		h1 = Math.imul(h3 ^ (h1 >>> 18), 597399067);
		h2 = Math.imul(h4 ^ (h2 >>> 22), 2869860233);
		h3 = Math.imul(h1 ^ (h3 >>> 17), 951274213);
		h4 = Math.imul(h2 ^ (h4 >>> 19), 2716044179);
		return (h1 ^ h2 ^ h3 ^ h4) >>> 0;
	}

	function mulberry32(a: number) {
		return function () {
			var t = (a += 0x6d2b79f5);
			t = Math.imul(t ^ (t >>> 15), t | 1);
			t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
			return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
		};
	}

	// --- 2. Color Logic (HSL) ---
	// Generates a palette where all colors are neighbors on the color wheel
	function generateMonochromaticPalette(rand: () => number) {
		const baseHue = Math.floor(rand() * 360);

		// Helper to keep hue within 0-360
		const wrap = (h: number) => (h < 0 ? h + 360 : h > 360 ? h - 360 : h);

		return {
			// Very dark, desaturated version of base hue for background
			bg: `hsl(${baseHue}, 20%, 8%)`,
			// Bright, punchy main color
			primary: {
				start: `hsl(${baseHue}, 75%, 60%)`,
				end: `hsl(${wrap(baseHue + 10)}, 85%, 50%)`
			},
			// Slightly shifted hue, lighter
			secondary: {
				start: `hsl(${wrap(baseHue - 25)}, 60%, 70%)`,
				end: `hsl(${wrap(baseHue - 15)}, 65%, 60%)`
			},
			// Opposing shift, darker or different saturation
			accent: {
				start: `hsl(${wrap(baseHue + 25)}, 70%, 75%)`,
				end: `hsl(${wrap(baseHue + 35)}, 80%, 65%)`
			}
		};
	}

	// --- 3. Spline Logic for Smooth Shapes ---
	// Creates smooth bezier curves between points instead of straight lines
	function getSplinePath(points: { x: number; y: number }[]) {
		if (points.length === 0) return '';

		// Duplicate first few points to close the loop smoothly
		const pts = [...points, points[0], points[1]];
		let d = `M ${(pts[0].x + pts[1].x) / 2} ${(pts[0].y + pts[1].y) / 2}`;

		for (let i = 1; i < pts.length - 1; i++) {
			const p1 = pts[i];
			const p2 = pts[i + 1];
			// Control point is p1, end point is midpoint between p1 and p2
			d += ` Q ${p1.x} ${p1.y} ${(p1.x + p2.x) / 2} ${(p1.y + p2.y) / 2}`;
		}
		return d + 'Z';
	}

	function generateOrganic(x: number, y: number, size: number, rand: () => number): string {
		const cx = x + size / 2;
		const cy = y + size / 2;
		const r = size / 2;
		const count = 7; // Number of points
		const points = [];

		for (let i = 0; i < count; i++) {
			const angle = (i / count) * Math.PI * 2;
			// Variance determines how "wobbly" the shape is
			const variance = 0.7 + rand() * 0.5;
			const radius = r * variance;
			points.push({
				x: cx + Math.cos(angle) * radius,
				y: cy + Math.sin(angle) * radius
			});
		}
		return getSplinePath(points);
	}

	// --- 4. State Derivation ---
	const config = $derived.by(() => {
		const seedNum = cyrb128(seed);
		const rand = mulberry32(seedNum);
		const palette = generateMonochromaticPalette(rand);

		const shapeTypes = ['circle', 'rect', 'organic'];

		const layers = [
			{
				id: 'l1',
				type: shapeTypes[Math.floor(rand() * 3)],
				x: 10 + rand() * 20,
				y: 10 + rand() * 20,
				size: 55 + rand() * 20,
				rotate: rand() * 360,
				gradient: palette.primary
			},
			{
				id: 'l2',
				type: shapeTypes[Math.floor(rand() * 3)],
				x: 35 + rand() * 30,
				y: 30 + rand() * 30,
				size: 40 + rand() * 20,
				rotate: rand() * 360,
				gradient: palette.secondary
			},
			{
				id: 'l3',
				type: shapeTypes[Math.floor(rand() * 3)],
				x: 50 + rand() * 20,
				y: 50 + rand() * 20,
				size: 25 + rand() * 15,
				rotate: rand() * 360,
				gradient: palette.accent
			}
		];

		return {
			palette,
			layers,
			uniqueId: `av-${seedNum}`
		};
	});
</script>

<svg
	viewBox="0 0 100 100"
	xmlns="http://www.w3.org/2000/svg"
	class={className}
	width={size}
	height={size}
	style="border-radius: 9999px; overflow: hidden; display: block; background-color: {config.palette
		.bg};"
>
	<defs>
		<filter id="{config.uniqueId}-paper">
			<feTurbulence type="fractalNoise" baseFrequency="0.8" numOctaves="3" result="noise" />
			<feDiffuseLighting in="noise" lighting-color="white" surfaceScale="1">
				<feDistantLight azimuth="45" elevation="60" />
			</feDiffuseLighting>
			<feComposite operator="in" in2="SourceGraphic" />
		</filter>

		<filter id="{config.uniqueId}-shadow" x="-50%" y="-50%" width="200%" height="200%">
			<feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#000" flood-opacity="0.25" />
			<feDropShadow dx="0" dy="1" stdDeviation="4" flood-color="#000" flood-opacity="0.1" />
		</filter>

		{#each config.layers as layer, i}
			<linearGradient id="{config.uniqueId}-g-{i}" x1="0%" y1="0%" x2="100%" y2="100%">
				<stop offset="0%" stop-color={layer.gradient.start} />
				<stop offset="100%" stop-color={layer.gradient.end} />
			</linearGradient>
		{/each}

		<mask id="{config.uniqueId}-mask">
			<circle cx="50" cy="50" r="50" fill="white" />
		</mask>
	</defs>

	<g mask="url(#{config.uniqueId}-mask)">
		{#each config.layers as layer, i (layer.id)}
			{@const seedNum = cyrb128(layer.id + seed)}
			{@const rand = mulberry32(seedNum)}

			<g transform={`rotate(${layer.rotate} 50 50)`} filter="url(#{config.uniqueId}-shadow)">
				{#if layer.type === 'circle'}
					<circle
						cx={layer.x + layer.size / 2}
						cy={layer.y + layer.size / 2}
						r={layer.size / 2}
						fill={`url(#${config.uniqueId}-g-${i})`}
					/>
				{:else if layer.type === 'rect'}
					<rect
						x={layer.x}
						y={layer.y}
						width={layer.size}
						height={layer.size}
						rx="8"
						fill={`url(#${config.uniqueId}-g-${i})`}
					/>
				{:else if layer.type === 'organic'}
					<path
						d={generateOrganic(layer.x, layer.y, layer.size, rand)}
						fill={`url(#${config.uniqueId}-g-${i})`}
					/>
				{/if}
			</g>
		{/each}

		<rect
			width="100"
			height="100"
			filter="url(#{config.uniqueId}-paper)"
			opacity="0.4"
			style="mix-blend-mode: soft-light;"
		/>

		<circle
			cx="50"
			cy="50"
			r="49"
			fill="none"
			stroke="white"
			stroke-opacity="0.1"
			stroke-width="1"
		/>
	</g>
</svg>
