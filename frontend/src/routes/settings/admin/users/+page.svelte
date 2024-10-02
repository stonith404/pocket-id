<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import UserService from '$lib/services/user-service';
	import appConfigStore from '$lib/stores/application-configuration-store';
	import type { Paginated } from '$lib/types/pagination.type';
	import type { User, UserCreate } from '$lib/types/user.type';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { LucideMinus } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import { slide } from 'svelte/transition';
	import UserForm from './user-form.svelte';
	import UserList from './user-list.svelte';

	let { data } = $props();
	let users: Paginated<User> = $state(data);
	let expandAddUser = $state(false);

	const userService = new UserService();

	async function createUser(user: UserCreate) {
		let success = true;
		await userService
			.create(user)
			.then(() => toast.success('User created successfully'))
			.catch((e) => {
				axiosErrorToast(e);
				success = false;
			});

		users = await userService.list();
		return success;
	}
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>

<Card.Root>
	<Card.Header>
		<div class="flex items-center justify-between">
			<div>
				<Card.Title>Create User</Card.Title>
				<Card.Description>Add a new user to {$appConfigStore.appName}.</Card.Description>
			</div>
			{#if !expandAddUser}
				<Button on:click={() => (expandAddUser = true)}>Add User</Button>
			{:else}
				<Button class="h-8 p-3" variant="ghost" on:click={() => (expandAddUser = false)}>
					<LucideMinus class="h-5 w-5" />
				</Button>
			{/if}
		</div>
	</Card.Header>
	{#if expandAddUser}
		<div transition:slide>
			<Card.Content>
				<UserForm callback={createUser} />
			</Card.Content>
		</div>
	{/if}
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Manage Users</Card.Title>
	</Card.Header>
	<Card.Content>
		<UserList {users} />
	</Card.Content>
</Card.Root>
