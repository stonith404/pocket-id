<script lang="ts">
	import { goto } from '$app/navigation';
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import Logo from '$lib/components/logo.svelte';
	import { Button } from '$lib/components/ui/button';
	import UserService from '$lib/services/user-service';
	import appConfigStore from '$lib/stores/application-configuration-store.js';
	import userStore from '$lib/stores/user-store.js';
	import type { User } from '$lib/types/user.type.js';
	import { axiosErrorToast } from '$lib/utils/error-util';

	let { data } = $props();
	let isLoading = $state(false);

	const userService = new UserService();

	async function authenticate() {
		isLoading = true;
		userService
			.exchangeOneTimeAccessToken(data.token)
			.then((user: User) => {
				userStore.setUser(user);
				goto('/settings');
			})
			.catch(axiosErrorToast);
		isLoading = false;
	}
</script>

<SignInWrapper>
	<div class="flex justify-center">
		<div class="bg-muted rounded-2xl p-3">
			<Logo class="h-10 w-10" />
		</div>
	</div>
	<h1 class="font-playfair mt-5 text-4xl font-bold">
		{data.token === 'setup' ? `${$appConfigStore.appName} Setup` : 'One Time Access'}
	</h1>
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
</SignInWrapper>
