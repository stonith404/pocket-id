<script lang="ts">
	import { goto } from '$app/navigation';
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import { Button } from '$lib/components/ui/button';
	import UserService from '$lib/services/user-service';
	import appConfigStore from '$lib/stores/application-configuration-store.js';
	import userStore from '$lib/stores/user-store.js';
	import { getAxiosErrorMessage } from '$lib/utils/error-util';
	import { onMount } from 'svelte';
	import LoginLogoErrorSuccessIndicator from '../components/login-logo-error-success-indicator.svelte';

	let { data } = $props();
	let isLoading = $state(false);
	let error: string | undefined = $state();
	const skipPage = data.redirect !== '/settings';

	const userService = new UserService();

	async function authenticate() {
		isLoading = true;
		try {
			const user = await userService.exchangeOneTimeAccessToken(data.token);
			userStore.setUser(user);

			try {
				goto(data.redirect);
			} catch (e) {
				error = 'Invalid redirect URL';
			}
		} catch (e) {
			error = getAxiosErrorMessage(e);
		}

		isLoading = false;
	}

	onMount(() => {
		if (skipPage) {
			authenticate();
		}
	});
</script>

<SignInWrapper>
	<div class="flex justify-center">
		<LoginLogoErrorSuccessIndicator error={!!error} />
	</div>
	<h1 class="font-playfair mt-5 text-4xl font-bold">
		{data.token === 'setup' ? `${$appConfigStore.appName} Setup` : 'One Time Access'}
	</h1>
	{#if error}
		<p class="text-muted-foreground mt-2">
			{error}. Please try again.
		</p>
	{:else if !skipPage}
		<p class="text-muted-foreground mt-2">
			{#if data.token === 'setup'}
				You're about to sign in to the initial admin account. Anyone with this link can access the
				account until a passkey is added. Please set up a passkey as soon as possible to prevent
				unauthorized access.
			{:else}
				You've been granted one-time access to your {$appConfigStore.appName} account. Please note that
				if you continue, this link will become invalid. To avoid this, make sure to add a passkey. Otherwise,
				you'll need to request a new link.
			{/if}
		</p>
		<Button class="mt-5" {isLoading} on:click={authenticate}>Continue</Button>
	{/if}
</SignInWrapper>
