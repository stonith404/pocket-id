<script lang="ts">
	import { goto } from '$app/navigation';
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import Logo from '$lib/components/logo.svelte';
	import { Button } from '$lib/components/ui/button';
	import UserService from '$lib/services/user-service';
	import applicationConfigurationStore from '$lib/stores/application-configuration-store.js';
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
			.then((user :User) => {
				userStore.setUser(user);
				goto('/settings')
			})
			.catch(axiosErrorToast);
		isLoading = false;
	}
</script>

<SignInWrapper>
	<div class="flex justify-center">
		<div class="rounded-2xl bg-muted p-3">
			<Logo class="h-10 w-10" />
		</div>
	</div>
	<h1 class="mt-5 font-playfair text-4xl font-bold">One Time Access</h1>
	<p class="mt-2 text-muted-foreground">
		You've been granted one-time access to your {$applicationConfigurationStore.appName} account. Please note that if you continue,
		this link will become invalid. To avoid this, make sure to add a passkey. Otherwise, you'll need
		to request a new link.
	</p>
	<Button class="mt-5" {isLoading} on:click={authenticate}>Continue</Button>
</SignInWrapper>
