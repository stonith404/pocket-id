<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import type { FormInput } from '$lib/utils/form-util';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { Input } from './ui/input';

	let {
		input = $bindable(),
		label,
		description,
		children,
		...restProps
	}: HTMLAttributes<HTMLDivElement> & {
		input?: FormInput<string | boolean | number>;
		label: string;
		description?: string;
		children?: Snippet;
	} = $props();

	const id = label.toLowerCase().replace(/ /g, '-');
</script>

<div {...restProps}>
	<Label class="mb-0" for={id}>{label}</Label>
	{#if description}
		<p class="text-muted-foreground mt-1 text-xs">{description}</p>
	{/if}
	<div class="mt-2">
		{#if children}
			{@render children()}
		{:else if input}
			<Input {id} bind:value={input.value} />
		{/if}
		{#if input?.error}
			<p class="mt-1 text-sm text-red-500">{input.error}</p>
		{/if}
	</div>
</div>
