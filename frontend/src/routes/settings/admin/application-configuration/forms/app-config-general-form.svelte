<script lang="ts">
	import FormInput from '$lib/components/form-input.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
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
		sessionDuration: appConfig.sessionDuration,
		emailsVerified: appConfig.emailsVerified,
		allowOwnAccountEdit: appConfig.allowOwnAccountEdit
	};

	const formSchema = z.object({
		appName: z.string().min(2).max(30),
		sessionDuration: z.number().min(1).max(43200),
		emailsVerified: z.boolean(),
		allowOwnAccountEdit: z.boolean()
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
			type="number"
			description="The duration of a session in minutes before the user has to sign in again."
			bind:input={$inputs.sessionDuration}
		/>
		<div class="items-top mt-5 flex space-x-2">
			<Checkbox id="admin-privileges" bind:checked={$inputs.allowOwnAccountEdit.value} />
			<div class="grid gap-1.5 leading-none">
				<Label for="admin-privileges" class="mb-0 text-sm font-medium leading-none">
					Enable Self-Account Editing
				</Label>
				<p class="text-muted-foreground text-[0.8rem]">
					Whether the user should be able to edit their own account details.
				</p>
			</div>
		</div>
		<div class="items-top mt-5 flex space-x-2">
			<Checkbox id="admin-privileges" bind:checked={$inputs.emailsVerified.value} />
			<div class="grid gap-1.5 leading-none">
				<Label for="admin-privileges" class="mb-0 text-sm font-medium leading-none">
					Emails Verified
				</Label>
				<p class="text-muted-foreground text-[0.8rem]">
					Whether the user's email should be marked as verified for the OIDC clients.
				</p>
			</div>
		</div>
	</div>
	<div class="mt-5 flex justify-end">
		<Button {isLoading} type="submit">Save</Button>
	</div>
</form>
