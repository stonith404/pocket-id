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
	let emailEnabled = $state(appConfig.emailEnabled == 'true');

	const updatedAppConfig = {
		emailEnabled: emailEnabled.toString(),
		smtpHost: appConfig.smtpHost,
		smtpPort: appConfig.smtpPort,
		smtpUser: appConfig.smtpUser,
		smtpPassword: appConfig.smtpPassword,
		smtpFrom: appConfig.smtpFrom
	};

	const formSchema = z.object({
		smtpHost: z.string().min(1),
		smtpPort: z.string().min(1),
		smtpUser: z.string().min(1),
		smtpPassword: z.string().min(1),
		smtpFrom: z.string().email()
	});

	const { inputs, ...form } = createForm< typeof formSchema>(formSchema, updatedAppConfig);

	async function onSubmit() {
		const data = form.validate();
		if (!data) return false;
		isLoading = true;
		await callback({
			...data,
			emailEnabled: 'true'
		}).finally(() => (isLoading = false));
		toast.success('Email configuration saved successfully');
		return true;
	}

	async function onDisable() {
		await callback({ emailEnabled: 'false' });
		emailEnabled = false;
		toast.success('Email disabled successfully');
	}

	async function onEnable() {
		if (await onSubmit()) {
			emailEnabled = true;
		}
	}
</script>

<form onsubmit={onSubmit}>
	<div class="mt-5 grid grid-cols-2 gap-5">
		<FormInput label="SMTP Host" bind:input={$inputs.smtpHost} />
		<FormInput label="SMTP Port" bind:input={$inputs.smtpPort} />
		<FormInput label="SMTP User" bind:input={$inputs.smtpUser} />
		<FormInput label="SMTP Password" type="password" bind:input={$inputs.smtpPassword} />
		<FormInput label="SMTP From" bind:input={$inputs.smtpFrom} />
	</div>
	<div class="mt-5 flex justify-end gap-3">
		{#if emailEnabled}
			<Button variant="secondary" onclick={onDisable}>Disable</Button>
			<Button {isLoading} onclick={onSubmit} type="submit">Save</Button>
		{:else}
			<Button {isLoading} onclick={onEnable} type="submit">Enable</Button>
		{/if}
	</div>
</form>
