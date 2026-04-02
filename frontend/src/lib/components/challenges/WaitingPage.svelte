<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Shield } from '@lucide/svelte';

	let { startTime, title = "Competition hasn't started yet" } = $props<{
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

	function updateCountdown() {
		if (!startTime) return;

		const start = new Date(startTime).getTime();
		const now = new Date().getTime();
		const diff = start - now;

		if (diff <= 0) {
			timeLeft = { days: 0, hours: 0, minutes: 0, seconds: 0, total: 0 };
			if (interval) clearInterval(interval);
			// Refresh page when it hits zero
			window.location.reload();
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

	onMount(() => {
		updateCountdown();
		interval = setInterval(updateCountdown, 1000);
	});

	onDestroy(() => {
		if (interval) clearInterval(interval);
	});

	const formattedStartTime = $derived(
		startTime
			? new Date(startTime).toLocaleString(undefined, {
					weekday: 'long',
					year: 'numeric',
					month: 'long',
					day: 'numeric',
					hour: '2-digit',
					minute: '2-digit'
				})
			: 'To be announced'
	);
</script>

<div class="flex flex-col items-center justify-center p-8 py-20 text-center">
	<div class="mb-8 space-y-4">
		<h1 class="text-4xl font-bold tracking-tight">
			{title}
		</h1>
		<p class="text-muted-foreground">
			The contest starts on {formattedStartTime}
		</p>
	</div>

	<!-- Simple Countdown -->
	<div class="grid w-full max-w-2xl grid-cols-2 gap-4 md:grid-cols-4">
		{#each [{ label: 'Days', value: timeLeft.days }, { label: 'Hours', value: timeLeft.hours }, { label: 'Minutes', value: timeLeft.minutes }, { label: 'Seconds', value: timeLeft.seconds }] as unit}
			<div class="bg-muted/50 flex flex-col items-center rounded-xl border p-4 md:p-6">
				<span class="font-mono text-3xl font-bold tabular-nums tracking-tighter md:text-5xl">
					{String(unit.value).padStart(2, '0')}
				</span>
				<span class="text-muted-foreground mt-1 text-[10px] font-bold uppercase tracking-wider"
					>{unit.label}</span
				>
			</div>
		{/each}
	</div>
</div>

<style>
</style>
