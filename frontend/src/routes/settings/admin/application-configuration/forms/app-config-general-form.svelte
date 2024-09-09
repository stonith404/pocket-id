<script lang="ts">
	import FormInput from '$lib/components/form-input.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { AllAppConfig } from '$lib/types/application-configuration';
	import { createForm } from '$lib/utils/form-util';
	import { toast } from 'svelte-sonner';
	import { z } from 'zod';

	let {
		callback,
		appConfig
	}: {
		appConfig: AllAppConfig;
		callback: (appConfig: Partial<AllAppConfig>) => Promise<void>;
	} = $props();

	let isLoading = $state(false);

	const updatedAppConfig = {
		appName: appConfig.appName,
		sessionDuration: appConfig.sessionDuration
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

	const { inputs, ...form } = createForm<typeof formSchema>(formSchema, updatedAppConfig);
	async function onSubmit() {
		const data = form.validate();
		if (!data) return;
		isLoading = true;
		await callback(data).finally(() => (isLoading = false));
		toast.success('Application configuration updated successfully');
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
