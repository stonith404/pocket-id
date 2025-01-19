<script lang="ts">
	import { browser } from '$app/environment';
	import { browserSupportsWebAuthn } from '@simplewebauthn/browser';
	import type { Snippet } from 'svelte';
	import { Button } from './ui/button';
	import * as Card from './ui/card';
	import WebAuthnUnsupported from './web-authn-unsupported.svelte';
	import { page } from '$app/stores';

	let {
		children,
		showEmailOneTimeAccessButton = false
	}: {
		children: Snippet;
		showEmailOneTimeAccessButton?: boolean;
	} = $props();
</script>

<!-- Desktop -->
<div class="hidden h-screen items-center text-center lg:flex">
	<div class="h-full min-w-[650px] p-16 {showEmailOneTimeAccessButton ? 'pb-0' : ''}">
		{#if browser && !browserSupportsWebAuthn()}
			<WebAuthnUnsupported />
		{:else}
			<div class="flex h-full flex-col">
				<div class="flex flex-grow flex-col items-center justify-center">
					{@render children()}
				</div>
				{#if showEmailOneTimeAccessButton}
					<div class="mb-4 flex justify-center">
						<Button
							href="/login/email?redirect={encodeURIComponent(
								$page.url.pathname + $page.url.search
							)}"
							variant="link"
							class="text-xs text-muted-foreground"
						>
							Don't have access to your passkey?
						</Button>
					</div>
				{/if}
			</div>
		{/if}
	</div>
	<img
		src="/api/application-configuration/background-image"
		class="h-screen w-[calc(100vw-650px)] rounded-l-[60px] object-cover"
		alt="Login background"
	/>
</div>

<!-- Mobile -->
<div
	class="flex h-screen items-center justify-center bg-[url('/api/application-configuration/background-image')] bg-cover bg-center text-center lg:hidden"
>
	<Card.Root class="mx-3">
		<Card.CardContent
			class="px-4 py-10 sm:p-10 {showEmailOneTimeAccessButton ? 'pb-3 sm:pb-3' : ''}"
		>
			{#if browser && !browserSupportsWebAuthn()}
				<WebAuthnUnsupported />
			{:else}
				{@render children()}
				{#if showEmailOneTimeAccessButton}
					<div class="mt-5">
						<Button
							href="/login/email?redirect={encodeURIComponent(
								$page.url.pathname + $page.url.search
							)}"
							variant="link"
							class="text-xs text-muted-foreground"
						>
							Don't have access to your passkey?
						</Button>
					</div>
				{/if}
			{/if}
		</Card.CardContent>
	</Card.Root>
</div>
