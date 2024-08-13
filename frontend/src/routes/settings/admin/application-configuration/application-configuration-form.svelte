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
		appName: applicationConfiguration.appName,
		sessionDuration: applicationConfiguration.sessionDuration
	};

	const formSchema = z.object({
		appName: z.string().min(2).max(30),
		sessionDuration: z.string().refine(
			(val) => {
				const num = Number(val);
				return Number.isInteger(num) && num >= 1 && num <= 43200;
			},
			{
				message: 'Session duration must be between 1 and 43200 minutes'
			}
		)
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
	<div class="flex flex-col gap-5">
		<FormInput label="Application Name" bind:input={$inputs.appName} />

		<FormInput
			label="Session Duration"
			description="The duration of a session in minutes before the user has to sign in again."
			bind:input={$inputs.sessionDuration}
		/>
	</div>
	<div class="mt-5 flex justify-end">
		<Button {isLoading} type="submit">Save</Button>
	</div>
</form>
