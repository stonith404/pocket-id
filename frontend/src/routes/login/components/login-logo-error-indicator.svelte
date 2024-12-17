<script lang="ts">
	import Logo from '$lib/components/logo.svelte';
	import CrossAnimated from '$lib/icons/cross-animated.svelte';
	import { fade } from 'svelte/transition';

	const {
		error
	}: {
		error: boolean;
	} = $props();
</script>

<div
	class="rounded-2xl p-3 transition-[background-color] duration-300
		 {error ? 'bg-red-200' : 'bg-muted'}"
>
	{#if error}
		<div class="flex h-10 w-10 items-center justify-center">
			<CrossAnimated class="h-5 w-5" />
		</div>
	{:else}
		<div in:fade={{ duration: 300 }}>
			<Logo class="h-10 w-10" />
		</div>
	{/if}
</div>
