<script lang="ts">
	import { Checkbox } from './ui/checkbox';
	import { Label } from './ui/label';

	let {
		id,
		checked = $bindable(),
		label,
		description,
		disabled = false,
		onCheckedChange
	}: {
		id: string;
		checked: boolean;
		label: string;
		description?: string;
		disabled?: boolean;
		onCheckedChange?: (checked: boolean) => void;
	} = $props();
</script>

<div class="items-top mt-5 flex space-x-2">
	<Checkbox
		{id}
		{disabled}
		onCheckedChange={(v) => onCheckedChange && onCheckedChange(v == true)}
		bind:checked
	/>
	<div class="grid gap-1.5 leading-none">
		<Label for={id} class="mb-0 text-sm font-medium leading-none">
			{label}
		</Label>
		{#if description}
			<p class="text-muted-foreground text-[0.8rem]">
				{description}
			</p>
		{/if}
	</div>
</div>
