<script lang="ts">
	import { page } from '$app/stores';
	import applicationConfigurationStore from '$lib/stores/application-configuration-store';
	import userStore from '$lib/stores/user-store';
	import Logo from '../logo.svelte';
	import HeaderAvatar from './header-avatar.svelte';

	let isAuthPage = $derived(
		!$page.error &&
			($page.url.pathname.startsWith('/authorize') || $page.url.pathname.startsWith('/login'))
	);
</script>

<div class=" w-full {isAuthPage ? 'absolute top-0 z-10 mt-4' : 'border-b'}">
	<div class="mx-auto flex w-full max-w-[1640px] items-center justify-between px-4 md:px-10">
		<div class="flex h-16 items-center">
			{#if !isAuthPage}
				<Logo class="mr-3 h-10 w-10" />
				<h1 class="text-lg font-medium" data-testid="application-name">
					{$applicationConfigurationStore.appName}
				</h1>
			{/if}
		</div>
		{#if $userStore?.id}
			<HeaderAvatar />
		{/if}
	</div>
</div>
