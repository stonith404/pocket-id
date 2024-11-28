<script lang="ts">
	import CheckboxWithLabel from '$lib/components/checkbox-with-label.svelte';
	import FormInput from '$lib/components/form-input.svelte';
	import { Button } from '$lib/components/ui/button';
	import AppConfigService from '$lib/services/app-config-service';
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

	const appConfigService = new AppConfigService();

	let isSendingTestEmail = $state(false);
	let emailEnabled = $state(appConfig.emailEnabled);

	const updatedAppConfig = {
		emailEnabled: appConfig.emailEnabled,
		smtpHost: appConfig.smtpHost,
		smtpPort: appConfig.smtpPort,
		smtpUser: appConfig.smtpUser,
		smtpPassword: appConfig.smtpPassword,
		smtpFrom: appConfig.smtpFrom,
		smtpTls: appConfig.smtpTls,
		smtpSkipCertVerify: appConfig.smtpSkipCertVerify
	};

	const formSchema = z.object({
		smtpHost: z.string().min(1),
		smtpPort: z.number().min(1),
		smtpUser: z.string(),
		smtpPassword: z.string(),
		smtpFrom: z.string().email(),
		smtpTls: z.boolean(),
		smtpSkipCertVerify: z.boolean()
	});

	const { inputs, ...form } = createForm<typeof formSchema>(formSchema, updatedAppConfig);

	async function onSubmit() {
		console.log('submit');
		const data = form.validate();
		if (!data) return false;
		await callback({
			...data,
			emailEnabled: true
		});
		toast.success('Email configuration updated successfully');
		return true;
	}

	async function onDisable() {
		emailEnabled = false;
		await callback({ emailEnabled });
		toast.success('Email disabled successfully');
	}

	async function onEnable() {
		if (await onSubmit()) {
			emailEnabled = true;
		}
	}

	async function onTestEmail() {
		isSendingTestEmail = true;
		await appConfigService
			.sendTestEmail()
			.then(() => toast.success('Test email sent successfully to your Email address.'))
			.catch(() =>
				toast.error('Failed to send test email. Check the server logs for more information.')
			)
			.finally(() => (isSendingTestEmail = false));
	}
</script>

<form onsubmit={onSubmit}>
	<div class="mt-5 grid grid-cols-1 gap-5 md:grid-cols-2 items-start">
		<FormInput label="SMTP Host" bind:input={$inputs.smtpHost} />
		<FormInput label="SMTP Port" type="number" bind:input={$inputs.smtpPort} />
		<FormInput label="SMTP User" bind:input={$inputs.smtpUser} />
		<FormInput label="SMTP Password" type="password" bind:input={$inputs.smtpPassword} />
		<FormInput label="SMTP From" bind:input={$inputs.smtpFrom} />
		<CheckboxWithLabel
			id="tls"
			label="TLS"
			description="Enable TLS for the SMTP connection."
			bind:checked={$inputs.smtpTls.value}
		/>
		<CheckboxWithLabel
			id="skip-cert-verify"
			label="Skip Certificate Verification"
			description="This can be useful for self-signed certificates."
			bind:checked={$inputs.smtpSkipCertVerify.value}
		/>
	</div>
	<div class="mt-8 flex flex-wrap justify-end gap-3">
		{#if emailEnabled}
			<Button variant="secondary" onclick={onDisable}>Disable</Button>
			<Button isLoading={isSendingTestEmail} variant="secondary" onclick={onTestEmail}
				>Send Test Email</Button
			>

			<Button onclick={onSubmit} type="submit">Save</Button>
		{:else}
			<Button onclick={onEnable}>Enable</Button>
		{/if}
	</div>
</form>
