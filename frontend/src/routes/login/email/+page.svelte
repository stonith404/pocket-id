<script lang="ts">
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import { Button } from '$lib/components/ui/button';
	import Input from '$lib/components/ui/input/input.svelte';
	import UserService from '$lib/services/user-service';
	import { fade } from 'svelte/transition';
	import LoginLogoErrorSuccessIndicator from '../components/login-logo-error-success-indicator.svelte';

	const { data } = $props();

	const userService = new UserService();

	let email = $state('');
	let isLoading = $state(false);
	let error: string | undefined = $state(undefined);
	let success = $state(false);

	async function requestEmail() {
		isLoading = true;
		await userService
			.requestOneTimeAccessEmail(email, data.redirect)
			.then(() => (success = true))
			.catch((e) => (error = e.response?.data.error || 'An unknown error occured'));

		isLoading = false;
	}
</script>

<svelte:head>
	<title>Email One Time Access</title>
</svelte:head>

<SignInWrapper>
	<div class="flex justify-center">
		<LoginLogoErrorSuccessIndicator {success} error={!!error} />
	</div>
	<h1 class="mt-5 font-playfair text-3xl font-bold sm:text-4xl">Email One Time Access</h1>
	{#if error}
		<p class="mt-2 text-muted-foreground" in:fade>
			{error}. Please try again.
		</p>
		<div class="mt-10 flex w-full justify-stretch gap-2">
			<Button variant="secondary" class="w-full" href="/">Go back</Button>
			<Button class="w-full" onclick={() => (error = undefined)}>Try again</Button>
		</div>
	{:else if success}
		<p class="mt-2 text-muted-foreground" in:fade>
			An email has been sent to the provided email, if it exists in the system.
		</p>
	{:else}
		<form onsubmit={requestEmail}>
			<p class="mt-2 text-muted-foreground" in:fade>
				Enter your email to receive an email with a one time access link.
			</p>
			<Input id="Email" class="mt-7" placeholder="Your email" bind:value={email} />
			<div class="mt-8 flex justify-stretch gap-2">
				<Button variant="secondary" class="w-full" href="/">Go back</Button>
				<Button class="w-full" type="submit" {isLoading}>Submit</Button>
			</div>
		</form>
	{/if}
</SignInWrapper>
