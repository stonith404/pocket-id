<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import UserService from '$lib/services/user-service';
	import WebAuthnService from '$lib/services/webauthn-service';
	import type { Passkey } from '$lib/types/passkey.type';
	import type { UserCreate } from '$lib/types/user.type';
	import { axiosErrorToast, getWebauthnErrorMessage } from '$lib/utils/error-util';
	import { startRegistration } from '@simplewebauthn/browser';
	import { toast } from 'svelte-sonner';
	import AccountForm from './account-form.svelte';
	import PasskeyList from './passkey-list.svelte';
	import RenamePasskeyModal from './rename-passkey-modal.svelte';

	let { data } = $props();
	let account = $state(data.account);
	let passkeys = $state(data.passkeys);
	let passkeyToRename: Passkey | null = $state(null);

	const userService = new UserService();
	const webauthnService = new WebAuthnService();

	async function updateAccount(user: UserCreate) {
		let success = true;
		await userService
			.updateCurrent(user)
			.then(() => toast.success('Account details updated successfully'))
			.catch((e) => {
				axiosErrorToast(e);
				success = false;
			});

		return success;
	}

	async function createPasskey() {
		try {
			const opts = await webauthnService.getRegistrationOptions();
			const attResp = await startRegistration(opts);
			const passkey = await webauthnService.finishRegistration(attResp);

			passkeys = await webauthnService.listCredentials();
			passkeyToRename = passkey;
		} catch (e) {
			toast.error(getWebauthnErrorMessage(e));
		}
	}
</script>

<svelte:head>
	<title>Account Settings</title>
</svelte:head>

<Card.Root>
	<Card.Header>
		<Card.Title>Account Details</Card.Title>
	</Card.Header>
	<Card.Content>
		<AccountForm {account} callback={updateAccount} />
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<div class="flex items-center justify-between">
			<div>
				<Card.Title>Passkeys</Card.Title>
				<Card.Description class="mt-1">
					Manage your passkeys that you can use to authenticate yourself.
				</Card.Description>
			</div>
			<Button size="sm" on:click={createPasskey}>Add Passkey</Button>
		</div>
	</Card.Header>
	{#if passkeys.length != 0}
		<Card.Content>
			<PasskeyList {passkeys} />
		</Card.Content>
	{/if}
</Card.Root>
<RenamePasskeyModal
	bind:passkey={passkeyToRename}
	callback={async () => (passkeys = await webauthnService.listCredentials())}
/>
