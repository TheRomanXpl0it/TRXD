<script lang="ts">
	import { Button } from '@/components/ui/button';
	import { Home, Search } from '@lucide/svelte';
</script>

<div class="flex min-h-[60vh] flex-col items-center justify-center text-center">
	<div class="mb-8 space-y-4">
		<h1 class="font-mono text-9xl font-bold text-gray-900 dark:text-white">404</h1>
		<p class="text-xl text-gray-600 dark:text-gray-400">Page not found</p>
		<p class="text-base text-gray-500 dark:text-gray-500">
			The page you're looking for doesn't exist or has been moved.
		</p>
	</div>

	<div class="flex flex-wrap justify-center gap-3">
		<a href="#/">
			<Button class="cursor-pointer">
				<Home class="mr-2 h-4 w-4" />
				Back to Home
			</Button>
		</a>
		<a href="#/challenges">
			<Button variant="outline" class="cursor-pointer">
				<Search class="mr-2 h-4 w-4" />
				View Challenges
			</Button>
		</a>
	</div>
</div>
