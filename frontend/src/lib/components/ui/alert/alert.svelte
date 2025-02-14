<script lang="ts">
	import { cn } from '$lib/utils/style.js';
	import { LucideX } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { type Variant, alertVariants } from './index.js';

	type $$Props = HTMLAttributes<HTMLDivElement> & {
		variant?: Variant;
		dismissibleId?: string;
	};

	let className: $$Props['class'] = undefined;
	export let variant: $$Props['variant'] = 'default';
	export let dismissibleId: $$Props['dismissibleId'] = undefined;
	export { className as class };

	let isVisible = !dismissibleId;

	onMount(() => {
		if (dismissibleId) {
			const dismissedAlerts = JSON.parse(localStorage.getItem('dismissed-alerts') || '[]');
			isVisible = !dismissedAlerts.includes(dismissibleId);
		}
	});

	function dismiss() {
		if (dismissibleId) {
			const dismissedAlerts = JSON.parse(localStorage.getItem('dismissed-alerts') || '[]');
			localStorage.setItem('dismissed-alerts', JSON.stringify([...dismissedAlerts, dismissibleId]));
			isVisible = false;
		}
	}
</script>

{#if isVisible}
	<div class={cn(alertVariants({ variant }), className)} {...$$restProps} role="alert">
		<slot />
		{#if dismissibleId}
			<button on:click={dismiss} class="absolute top-0 right-0 m-3 text-black dark:text-white"><LucideX class="w-4" /></button>
		{/if}
	</div>
{/if}
