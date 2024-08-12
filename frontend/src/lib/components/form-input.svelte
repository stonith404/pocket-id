<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import type { FormInput } from '$lib/utils/form-util';
	import type { Snippet } from 'svelte';
	import { Input } from './ui/input';

	let {
		input = $bindable(),
		label,
		children
	}: {
		input: FormInput<string | boolean | number>;
		label: string;
		children?: Snippet;
	} = $props();

	const id = label.toLowerCase().replace(/ /g, '-');
</script>

<div>
	<Label for={id}>{label}</Label>
	{#if children}
		{@render children()}
	{:else}
		<Input {id} bind:value={input.value} />
	{/if}
	{#if input.error}
		<p class="text-sm text-red-500">{input.error}</p>
	{/if}
</div>
