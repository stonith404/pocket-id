<script lang="ts">
	import FileInput from '$lib/components/file-input.svelte';
	import FormInput from '$lib/components/form-input.svelte';
	import { Button } from '$lib/components/ui/button';
	import Label from '$lib/components/ui/label/label.svelte';
	import type {
		OidcClient,
		OidcClientCreate,
		OidcClientCreateWithLogo
	} from '$lib/types/oidc.type';
	import { createForm } from '$lib/utils/form-util';
	import { z } from 'zod';

	let {
		callback,
		existingClient
	}: {
		existingClient?: OidcClient;
		callback: (user: OidcClientCreateWithLogo) => Promise<boolean>;
	} = $props();

	let isLoading = $state(false);
	let logo = $state<File | null>(null);
	let logoDataURL: string | null = $state(
		existingClient?.hasLogo ? `/api/oidc/clients/${existingClient!.id}/logo` : null
	);

	const client: OidcClientCreate = {
		name: existingClient?.name || '',
		callbackURL: existingClient?.callbackURL || ''
	};

	const formSchema = z.object({
		name: z.string().min(2).max(50),
		callbackURL: z.string().url()
	});

	type FormSchema = typeof formSchema;
	const { inputs, ...form } = createForm<FormSchema>(formSchema, client);

	async function onSubmit() {
		const data = form.validate();
		if (!data) return;
		isLoading = true;
		const success = await callback({
			...data,
			logo
		});
		// Reset form if client was successfully created
		if (success && !existingClient) form.reset();
		isLoading = false;
	}

	function onLogoChange(e: Event) {
		const file = (e.target as HTMLInputElement).files?.[0] || null;
		if (file) {
			logo = file;
			const reader = new FileReader();
			reader.onload = (event) => {
				logoDataURL = event.target?.result as string;
			};
			reader.readAsDataURL(file);
		}
	}

	function resetLogo() {
		logo = null;
		logoDataURL = null;
	}
</script>

<form onsubmit={onSubmit}>
	<div class="mt-3 grid grid-cols-2 gap-3">
		<FormInput label="Name" bind:input={$inputs.name} />
		<FormInput label="Callback URL" bind:input={$inputs.callbackURL} />
		<div class="mt-3">
			<Label for="logo">Logo</Label>
			<div class="mt-2 flex items-end gap-3">
				{#if logoDataURL}
					<div class="h-32 w-32 rounded-2xl bg-muted p-3">
						<img class="m-auto max-h-full max-w-full object-contain" src={logoDataURL} alt={`${$inputs.name.value} logo`} />
					</div>
				{/if}
				<div class="flex flex-col gap-2">
					<FileInput
						id="logo"
						variant="secondary"
						accept="image/png, image/jpeg, image/svg+xml"
						onchange={onLogoChange}
					>
						<Button variant="secondary">
							{existingClient?.hasLogo ? 'Change Logo' : 'Upload Logo'}
						</Button>
					</FileInput>
					{#if logoDataURL}
						<Button variant="outline" on:click={resetLogo}>Remove Logo</Button>
					{/if}
				</div>
			</div>
		</div>
	</div>
	<div class="w-full"></div>
	<div class="mt-5 flex justify-end">
		<Button {isLoading} type="submit">Save</Button>
	</div>
</form>
