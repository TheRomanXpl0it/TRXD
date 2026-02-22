<script lang="ts">
	let {
		children,
		variant = 'default',
		class: className = ''
	} = $props<{
		children: any;
		variant?: 'default' | 'category' | 'type';
		class?: string;
	}>();

	const baseClasses = 'inline-flex items-center rounded px-2 py-0.5 text-xs font-medium';

	const variantClasses = {
		default: 'border dark:border-gray-700',
		category: 'border dark:border-gray-700',
		type: 'rounded-md bg-gray-100 dark:bg-gray-700 dark:text-gray-300 text-gray-600'
	};
</script>

<span
	class="{baseClasses} {variantClasses[variant as 'default' | 'category' | 'type']} {className}"
>
	{@render children()}
</span>
