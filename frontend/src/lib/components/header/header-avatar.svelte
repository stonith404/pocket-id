<script lang="ts">
	import { goto } from '$app/navigation';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import WebAuthnService from '$lib/services/webauthn-service';
	import userStore from '$lib/stores/user-store';
	import { LucideLogOut, LucideUser } from 'lucide-svelte';

	const webauthnService = new WebAuthnService();

	let initials = $derived(
		($userStore!.firstName.charAt(0) + $userStore!.lastName?.charAt(0)).toUpperCase()
	);

	function logout() {
		webauthnService.logout();
		window.location.reload();
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		><Avatar.Root>
			<Avatar.Fallback>{initials}</Avatar.Fallback>
		</Avatar.Root></DropdownMenu.Trigger
	>
	<DropdownMenu.Content class="w-40" align="start">
		<DropdownMenu.Label class="font-normal">
			<div class="flex flex-col space-y-1">
				<p class="text-sm font-medium leading-none">
					{$userStore?.firstName}
					{$userStore?.lastName}
				</p>
				<p class="text-xs leading-none text-muted-foreground">{$userStore?.email}</p>
			</div>
		</DropdownMenu.Label>
		<DropdownMenu.Separator />
		<DropdownMenu.Group>
			<DropdownMenu.Item href="/settings/account"
				><LucideUser class="mr-2 h-4 w-4" /> My Account</DropdownMenu.Item
			>
			<DropdownMenu.Item on:click={logout}
				><LucideLogOut class="mr-2 h-4 w-4" /> Logout</DropdownMenu.Item
			>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
