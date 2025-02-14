<script lang="ts">
	import { goto } from '$app/navigation';
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import Logo from '$lib/components/logo.svelte';
	import { Button } from '$lib/components/ui/button';
	import WebAuthnService from '$lib/services/webauthn-service';
	import userStore from '$lib/stores/user-store.js';
	import { axiosErrorToast } from '$lib/utils/error-util.js';

	let isLoading = $state(false);

	const webauthnService = new WebAuthnService();

	async function signOut() {
		isLoading = true;
		await webauthnService
			.logout()
			.then(() => goto('/'))
			.catch(axiosErrorToast);
		isLoading = false;
	}
</script>

<svelte:head>
	<title>Logout</title>
</svelte:head>

<SignInWrapper>
	<div class="flex justify-center">
		<div class="bg-muted rounded-2xl p-3">
			<Logo class="h-10 w-10" />
		</div>
	</div>
	<h1 class="font-playfair mt-5 text-4xl font-bold">Sign out</h1>

	<p class="text-muted-foreground mt-2">
		Do you want to sign out of Pocket ID with the account <b>{$userStore?.username}</b>?
	</p>
	<div class="mt-10 flex w-full justify-stretch gap-2">
		<Button class="w-full" variant="secondary" onclick={() => history.back()}>Cancel</Button>
		<Button class="w-full" {isLoading} onclick={signOut}>Sign out</Button>
	</div>
</SignInWrapper>
