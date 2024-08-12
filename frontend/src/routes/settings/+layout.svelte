<script lang="ts">
	import { page } from '$app/stores';
	import userStore from '$lib/stores/user-store';
	import type { Snippet } from 'svelte';

	let {
		children
	}: {
		children: Snippet;
	} = $props();

	let links = $state([{ href: '/settings/account', label: 'My Account' }]);

	if ($userStore?.isAdmin) {
		links = [
			...links,
			{ href: '/settings/admin/users', label: 'Users' },
			{ href: '/settings/admin/oidc-clients', label: 'OIDC Clients' },
			{ href: '/settings/admin/application-configuration', label: 'Application Configuration' }
		];
	}
</script>

<section>
	<div class="h-screen w-full">
		<main class="flex min-h-screen flex-1 flex-col gap-4 bg-muted/40 p-4 md:gap-8 md:p-10">
			<div class="mx-auto grid w-full max-w-[1440px] gap-2">
				<h1 class="text-3xl font-semibold">Settings</h1>
			</div>
			<div
				class="mx-auto grid w-full max-w-[1440px] items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
			>
				<nav class="grid gap-4 text-sm text-muted-foreground">
					{#each links as { href, label }}
						<a {href} class={$page.url.pathname.startsWith(href) ? 'font-bold text-primary' : ''}>
							{label}
						</a>
					{/each}
				</nav>
				<div class="flex flex-col gap-5">
					{@render children()}
				</div>
			</div>
		</main>
	</div>
</section>
