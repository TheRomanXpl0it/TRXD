<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Code, Database, Lock, Shield, Coffee, Command, Cpu, Radio } from '@lucide/svelte';

	let { startTime, title = 'Starting soon' } = $props<{
		startTime: string | null;
		title?: string;
	}>();

	let timeLeft = $state({
		days: 0,
		hours: 0,
		minutes: 0,
		seconds: 0,
		total: 0
	});

	let interval: any;
	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;

	function updateCountdown() {
		if (!startTime) return;

		const start = new Date(startTime).getTime();
		const now = new Date().getTime();
		const diff = start - now;

		if (diff <= 0) {
			timeLeft = { days: 0, hours: 0, minutes: 0, seconds: 0, total: 0 };
			if (interval) clearInterval(interval);
			if (typeof window !== 'undefined') window.location.reload();
			return;
		}

		timeLeft = {
			days: Math.floor(diff / (1000 * 60 * 60 * 24)),
			hours: Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
			minutes: Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60)),
			seconds: Math.floor((diff % (1000 * 60)) / 1000),
			total: diff
		};
	}

	// Matrix Rain Logic
	let animationFrame: number;
	let lastTime = 0;
	const fps = 20; // Throttled FPS for slower movement
	const characters = '0123456789ABCDEFHIJKLMNOPQRSTUVWXYZ@#$%^&*()*&^%';
	let columns = 0;
	let drops: number[] = [];

	function initMatrix() {
		if (!canvas) return;
		ctx = canvas.getContext('2d');
		resizeCanvas();
		window.addEventListener('resize', resizeCanvas);
	}

	function resizeCanvas() {
		if (!canvas || !canvas.parentElement) return;
		canvas.width = canvas.parentElement.clientWidth;
		canvas.height = canvas.parentElement.clientHeight;
		columns = Math.floor(canvas.width / 20);
		// Randomize starting positions so they don't drop all at once
		drops = Array(columns)
			.fill(0)
			.map(() => Math.floor(Math.random() * -100));
	}

	function drawMatrix(currentTime = 0) {
		if (!ctx || !canvas) return;

		animationFrame = requestAnimationFrame(drawMatrix);

		const delta = currentTime - lastTime;
		if (delta < 1000 / fps) return;

		lastTime = currentTime;

		const isDark = document.documentElement.classList.contains('dark');

		// Semitransparent background to create trail effect
		// In dark mode use black, in light mode use the background color (usually white)
		ctx.fillStyle = isDark ? 'rgba(0, 0, 0, 0.15)' : 'rgba(255, 255, 255, 0.15)';
		ctx.fillRect(0, 0, canvas.width, canvas.height);

		// Matrix characters
		// Grey in dark mode, darker grey in light mode for contrast
		ctx.fillStyle = isDark ? 'rgba(120, 120, 120, 0.4)' : 'rgba(100, 100, 100, 0.25)';
		ctx.font = '16px monospace';

		for (let i = 0; i < drops.length; i++) {
			const text = characters.charAt(Math.floor(Math.random() * characters.length));
			ctx.fillText(text, i * 20, drops[i] * 20);

			if (drops[i] * 20 > canvas.height && Math.random() > 0.975) {
				drops[i] = 0;
			}
			drops[i]++;
		}
	}

	onMount(() => {
		updateCountdown();
		interval = setInterval(updateCountdown, 1000);
		initMatrix();
		drawMatrix();
	});

	onDestroy(() => {
		if (interval) clearInterval(interval);
		if (typeof window !== 'undefined') {
			window.removeEventListener('resize', resizeCanvas);
			cancelAnimationFrame(animationFrame);
		}
	});

	const formattedStartTime = $derived(
		startTime
			? new Date(startTime).toLocaleString('en-GB', {
					weekday: 'long',
					year: 'numeric',
					month: 'long',
					day: 'numeric',
					hour: '2-digit',
					minute: '2-digit'
				})
			: 'To be announced'
	);

	const floatingIcons = [
		{ icon: Code, top: '15%', left: '10%', delay: '0s', size: 24, speed: '12s' },
		{ icon: Database, top: '25%', left: '85%', delay: '1.5s', size: 32, speed: '15s' },
		{ icon: Lock, top: '65%', left: '5%', delay: '0.8s', size: 28, speed: '10s' },
		{ icon: Shield, top: '75%', left: '90%', delay: '2.2s', size: 20, speed: '18s' }
	];
</script>

<div
	class="border-border bg-background text-foreground relative flex min-h-[85vh] flex-col items-center justify-center overflow-hidden rounded-3xl border p-8 py-24 text-center mt-8"
>
	<!-- Matrix Background -->
	<canvas bind:this={canvas} class="absolute inset-0 z-0 h-full w-full"></canvas>

	<!-- Vignette Effect (Reactive to theme) -->
	<div
		class="z-1 bg-radial-gradient to-background/80 pointer-events-none absolute inset-0 from-transparent"
	></div>

	<!-- Floating Elements (Integrated with Matrix Theme) -->
	{#each floatingIcons as item}
		<div
			class="animate-float z-2 text-primary/40 pointer-events-none absolute"
			style="top: {item.top}; left: {item.left}; --float-delay: {item.delay}; --float-speed: {item.speed};"
		>
			<item.icon size={item.size} strokeWidth={1} />
		</div>
	{/each}

	<!-- Foreground Content -->
	<div
		class="glass-container border-border bg-card/60 relative z-10 w-full max-w-4xl space-y-12 rounded-[2rem] border p-8 shadow-2xl backdrop-blur-3xl md:p-16"
	>
		<div class="space-y-4 pb-6">
			<h1 class="text-foreground pb-2 text-4xl font-black leading-[1.2] tracking-tight md:text-6xl">
				{title}
			</h1>
			<p class="mx-auto max-w-xl text-lg leading-relaxed opacity-90">
				The CTF starts on <span class="font-bold">{formattedStartTime}</span>.
			</p>
			<p class="mx-auto max-w-xl text-lg leading-relaxed">Get your horses ready.</p>
		</div>

		<!-- Glass Countdown -->
		<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
			{#each [{ label: 'Days', value: timeLeft.days }, { label: 'Hours', value: timeLeft.hours }, { label: 'Minutes', value: timeLeft.minutes }, { label: 'Seconds', value: timeLeft.seconds }] as unit}
				<div
					class="border-border bg-muted/20 group flex flex-col items-center rounded-2xl border p-6"
				>
					<span
						class="text-primary font-mono text-4xl font-black tabular-nums tracking-tighter md:text-6xl"
					>
						{String(unit.value).padStart(2, '0')}
					</span>
					<span
						class="text-muted-foreground mt-2 text-[10px] font-black uppercase tracking-[0.2em]"
					>
						{unit.label}
					</span>
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	:root {
		--primary-rgb: 59, 130, 246;
	}

	.bg-radial-gradient {
		background: radial-gradient(circle at center, transparent 0%, var(--background) 100%);
	}

	.glass-container {
		box-shadow: 0 0 40px rgba(0, 0, 0, 0.5);
	}

	@keyframes float {
		0% {
			transform: translate(0, 0) rotate(0deg);
		}
		33% {
			transform: translate(10px, -15px) rotate(3deg);
		}
		66% {
			transform: translate(-5px, -25px) rotate(-3deg);
		}
		100% {
			transform: translate(0, 0) rotate(0deg);
		}
	}

	.animate-float {
		animation: float var(--float-speed) ease-in-out infinite;
		animation-delay: var(--float-delay);
	}
</style>
