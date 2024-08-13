<script lang="ts">
	import { browser } from '$app/environment';
	import { browserSupportsWebAuthn } from '@simplewebauthn/browser';
	import type { Snippet } from 'svelte';
	import * as Card from './ui/card';
	import WebAuthnUnsupported from './web-authn-unsupported.svelte';

	let {
		children
	}: {
		children: Snippet;
	} = $props();
</script>

<div class="hidden h-screen items-center text-center lg:flex">
	<div class="min-w-[650px] p-16">
		{#if browser && !browserSupportsWebAuthn()}
			<WebAuthnUnsupported />
		{:else}
			{@render children()}
		{/if}
	</div>
	<img
		src="/api/application-configuration/background-image"
		class="h-screen w-[calc(100vw-650px)] rounded-l-[60px] object-cover"
		alt="Login background"
	/>
</div>

<div
	class="flex h-screen items-center justify-center bg-[url('/api/application-configuration/background-image')] bg-cover bg-center text-center lg:hidden"
>
	<Card.Root class="mx-3">
		<Card.CardContent class="px-4 py-10 sm:p-10">
			{#if browser && !browserSupportsWebAuthn()}
				<WebAuthnUnsupported />
			{:else}
				{@render children()}
			{/if}
		</Card.CardContent>
	</Card.Root>
</div>
