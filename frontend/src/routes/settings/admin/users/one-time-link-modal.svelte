<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';

	let {
		oneTimeLink = $bindable()
	}: {
		oneTimeLink: string | null;
	} = $props();

	function onOpenChange(open: boolean) {
		if (!open) {
			oneTimeLink = null;
		}
	}
</script>

<Dialog.Root open={!!oneTimeLink} {onOpenChange}>
	<Dialog.Content class="max-w-md">
		<Dialog.Header>
			<Dialog.Title>One Time Link</Dialog.Title>
			<Dialog.Description
				>Use this link to sign in once. This is needed for users who haven't added a passkey yet or
				have lost it.</Dialog.Description
			>
		</Dialog.Header>
		<div>
			<Label for="one-time-link">One Time Link</Label>
			<Input id="one-time-link" value={oneTimeLink} readonly />
		</div>
	</Dialog.Content>
</Dialog.Root>
