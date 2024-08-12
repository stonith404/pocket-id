<script lang="ts">
	import FormInput from '$lib/components/form-input.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { AllApplicationConfiguration } from '$lib/types/application-configuration';
	import { createForm } from '$lib/utils/form-util';
	import { z } from 'zod';

	let {
		callback,
		applicationConfiguration
	}: {
		applicationConfiguration: AllApplicationConfiguration;
		callback: (user: AllApplicationConfiguration) => Promise<void>;
	} = $props();

	let isLoading = $state(false);

	const updatedApplicationConfiguration: AllApplicationConfiguration = {
		appName: applicationConfiguration.appName
	};

	const formSchema = z.object({
		appName: z.string().min(2).max(30)
	});
	type FormSchema = typeof formSchema;

	const { inputs, ...form } = createForm<FormSchema>(formSchema, updatedApplicationConfiguration);
	async function onSubmit() {
		const data = form.validate();
		if (!data) return;
		isLoading = true;
		await callback(data);
		isLoading = false;
	}
</script>

<form onsubmit={onSubmit}>
	<div class="flex gap-3">
		<div class="w-full">
			<FormInput label="Application Name" bind:input={$inputs.appName} />
		</div>
	</div>
	<div class="mt-5 flex justify-end">
		<Button {isLoading} type="submit">Save</Button>
	</div>
</form>
