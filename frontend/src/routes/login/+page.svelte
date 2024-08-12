<script>
	import { goto } from '$app/navigation';
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import Logo from '$lib/components/logo.svelte';
	import { Button } from '$lib/components/ui/button';
	import WebAuthnService from '$lib/services/webauthn-service';
	import applicationConfigurationStore from '$lib/stores/application-configuration-store';
	import userStore from '$lib/stores/user-store';
	import { getWebauthnErrorMessage } from '$lib/utils/error-util';
	import { startAuthentication } from '@simplewebauthn/browser';
	import { toast } from 'svelte-sonner';
	const webauthnService = new WebAuthnService();

	let isLoading = $state(false);

	async function authenticate() {
		isLoading = true;
		try {
			const loginOptions = await webauthnService.getLoginOptions();
			const authResponse = await startAuthentication(loginOptions);
			const user = await webauthnService.finishLogin(authResponse);

			userStore.setUser(user);
			goto('/settings');
		} catch (e) {
			toast.error(getWebauthnErrorMessage(e));
		}
		isLoading = false;
	}
</script>

<svelte:head>
	<title>Sign In</title>
</svelte:head>

<SignInWrapper>
	<div class="flex justify-center">
		<div class="bg-muted rounded-2xl p-3">
			<Logo class="h-10 w-10" />
		</div>
	</div>
	<h1 class="font-playfair mt-5 text-3xl font-bold sm:text-4xl">
		Sign in to {$applicationConfigurationStore.appName}
	</h1>
	<p class="text-muted-foreground mt-2">
		Authenticate yourself with your passkey to access the admin panel
	</p>
	<Button class="mt-5" {isLoading} on:click={authenticate}>Authenticate</Button>
</SignInWrapper>
