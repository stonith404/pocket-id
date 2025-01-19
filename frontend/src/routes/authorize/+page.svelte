<script lang="ts">
	import SignInWrapper from '$lib/components/login-wrapper.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import OidcService from '$lib/services/oidc-service';
	import WebAuthnService from '$lib/services/webauthn-service';
	import appConfigStore from '$lib/stores/application-configuration-store';
	import userStore from '$lib/stores/user-store';
	import { getWebauthnErrorMessage } from '$lib/utils/error-util';
	import { startAuthentication } from '@simplewebauthn/browser';
	import { AxiosError } from 'axios';
	import { LucideMail, LucideUser, LucideUsers } from 'lucide-svelte';
	import { slide } from 'svelte/transition';
	import type { PageData } from './$types';
	import ClientProviderImages from './components/client-provider-images.svelte';
	import ScopeItem from './components/scope-item.svelte';

	const webauthnService = new WebAuthnService();
	const oidService = new OidcService();

	let isLoading = false;
	let success = false;
	let errorMessage: string | null = null;
	let authorizationRequired = false;

	export let data: PageData;
	let { scope, nonce, client, state, callbackURL, codeChallenge, codeChallengeMethod } = data;

	async function authorize() {
		isLoading = true;
		try {
			// Get access token if not signed in
			if (!$userStore?.id) {
				const loginOptions = await webauthnService.getLoginOptions();
				const authResponse = await startAuthentication(loginOptions);
				await webauthnService.finishLogin(authResponse);
			}

			await oidService
				.authorize(client!.id, scope, callbackURL, nonce, codeChallenge, codeChallengeMethod)
				.then(async ({ code, callbackURL }) => {
					onSuccess(code, callbackURL);
				});
		} catch (e) {
			if (e instanceof AxiosError && e.response?.status === 403) {
				authorizationRequired = true;
			} else {
				errorMessage = getWebauthnErrorMessage(e);
			}
			isLoading = false;
		}
	}

	async function authorizeNewClient() {
		isLoading = true;
		try {
			await oidService
				.authorizeNewClient(
					client!.id,
					scope,
					callbackURL,
					nonce,
					codeChallenge,
					codeChallengeMethod
				)
				.then(async ({ code, callbackURL }) => {
					onSuccess(code, callbackURL);
				});
		} catch (e) {
			errorMessage = getWebauthnErrorMessage(e);
			isLoading = false;
		}
	}

	function onSuccess(code: string, callbackURL: string) {
		success = true;
		setTimeout(() => {
			const redirectURL = new URL(callbackURL);
			redirectURL.searchParams.append('code', code);
			redirectURL.searchParams.append('state', state);

			window.location.href = redirectURL.toString();
		}, 1000);
	}
</script>

<svelte:head>
	<title>Sign in to {client.name}</title>
</svelte:head>

{#if client == null}
	<p>Client not found</p>
{:else}
	<SignInWrapper showEmailOneTimeAccessButton={$appConfigStore.emailOneTimeAccessEnabled}>
		<ClientProviderImages {client} {success} error={!!errorMessage} />
		<h1 class="font-playfair mt-5 text-3xl font-bold sm:text-4xl">Sign in to {client.name}</h1>
		{#if errorMessage}
			<p class="text-muted-foreground mb-10 mt-2">
				{errorMessage}. Please try again.
			</p>
		{/if}
		{#if !authorizationRequired && !errorMessage}
			<p class="text-muted-foreground mb-10 mt-2">
				Do you want to sign in to <b>{client.name}</b> with your
				<b>{$appConfigStore.appName}</b> account?
			</p>
		{:else if authorizationRequired}
			<div transition:slide={{ duration: 300 }}>
				<Card.Root class="mb-10 mt-6">
					<Card.Header class="pb-5">
						<p class="text-muted-foreground text-start">
							<b>{client.name}</b> wants to access the following information:
						</p>
					</Card.Header>
					<Card.Content data-testid="scopes">
						<div class="flex flex-col gap-3">
							{#if scope!.includes('email')}
								<ScopeItem icon={LucideMail} name="Email" description="View your email address" />
							{/if}
							{#if scope!.includes('profile')}
								<ScopeItem
									icon={LucideUser}
									name="Profile"
									description="View your profile information"
								/>
							{/if}
							{#if scope!.includes('groups')}
								<ScopeItem
									icon={LucideUsers}
									name="Groups"
									description="View the groups you are a member of"
								/>
							{/if}
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		{/if}
		<div class="flex w-full justify-stretch gap-2">
			<Button onclick={() => history.back()} class="w-full" variant="secondary">Cancel</Button>
			{#if !errorMessage}
				<Button
					class="w-full"
					{isLoading}
					on:click={authorizationRequired ? authorizeNewClient : authorize}
				>
					Sign in
				</Button>
			{:else}
				<Button class="w-full" on:click={() => (errorMessage = null)}>Try again</Button>
			{/if}
		</div>
	</SignInWrapper>
{/if}
