<script lang="ts">
	import Logo from '$lib/components/logo.svelte';
	import CheckmarkAnimated from '$lib/icons/checkmark-animated.svelte';
	import CrossAnimated from '$lib/icons/cross-animated.svelte';
	import { fade } from 'svelte/transition';

	const {
		error,
		success
	}: {
		error?: boolean;
		success?: boolean;
	} = $props();
</script>

<div
	class="rounded-2xl p-3 transition-[background-color] duration-300
		 {error ? 'bg-red-200' : success ? 'bg-green-200' : 'bg-muted'}"
>
	{#if error || success}
		<div class="flex h-10 w-10 items-center justify-center">
			{#if error}
				<CrossAnimated class="h-5 w-5" />
			{:else}
				<CheckmarkAnimated class="h-5 w-5" />
			{/if}
		</div>
	{:else}
		<div in:fade={{ duration: 300 }}>
			<Logo class="h-10 w-10" />
		</div>
	{/if}
</div>
