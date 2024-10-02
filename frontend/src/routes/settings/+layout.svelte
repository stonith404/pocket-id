<script lang="ts">
	import { page } from '$app/stores';
	import userStore from '$lib/stores/user-store';
	import type { Snippet } from 'svelte';

	let {
		children
	}: {
		children: Snippet;
	} = $props();

	let links = $state([
		{ href: '/settings/account', label: 'My Account' },
		{ href: '/settings/audit-log', label: 'Audit Log' }
	]);

	if ($userStore?.isAdmin) {
		links = [
			...links,
			{ href: '/settings/admin/users', label: 'Users' },
			{ href: '/settings/admin/user-groups', label: 'User Groups' },
			{ href: '/settings/admin/oidc-clients', label: 'OIDC Clients' },
			{ href: '/settings/admin/application-configuration', label: 'Application Configuration' }
		];
	}
</script>

<section>
	<div class="bg-muted/40 min-h-screen w-full">
		<main class="mx-auto flex max-w-[1640px] flex-col gap-x-4 gap-y-10 p-4 md:p-10 lg:flex-row">
			<div>
				<div class="mx-auto grid w-full gap-2">
					<h1 class="mb-5 text-3xl font-semibold">Settings</h1>
				</div>
				<div
					class="mx-auto grid items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
				>
					<nav class="text-muted-foreground grid gap-4 text-sm">
						{#each links as { href, label }}
							<a {href} class={$page.url.pathname.startsWith(href) ? 'text-primary font-bold' : ''}>
								{label}
							</a>
						{/each}
					</nav>
				</div>
			</div>
			<div class="flex w-full flex-col gap-5">
				{@render children()}
			</div>
		</main>
	</div>
</section>
