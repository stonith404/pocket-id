<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import UserGroupService from '$lib/services/user-group-service';
	import UserService from '$lib/services/user-service';
	import type { UserGroupCreate } from '$lib/types/user-group.type';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { LucideChevronLeft } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import UserGroupForm from '../user-group-form.svelte';
	import UserSelection from '../user-selection.svelte';

	let { data } = $props();
	let userGroup = $state({
		...data.userGroup,
		userIds: data.userGroup.users.map((u) => u.id)
	});

	const userGroupService = new UserGroupService();
	const userService = new UserService();

	async function updateUserGroup(updatedUserGroup: UserGroupCreate) {
		let success = true;
		await userGroupService
			.update(userGroup.id, updatedUserGroup)
			.then(() => toast.success('User Group updated successfully'))
			.catch((e) => {
				axiosErrorToast(e);
				success = false;
			});

		return success;
	}

	async function updateUserGroupUsers(userIds: string[]) {
		await userGroupService
			.updateUsers(userGroup.id, userIds)
			.then(() => toast.success('Users updated successfully'))
			.catch((e) => {
				axiosErrorToast(e);
			});
	}
</script>

<svelte:head>
	<title>User Group Details {userGroup.name}</title>
</svelte:head>

<div>
	<a class="text-muted-foreground flex text-sm" href="/settings/admin/user-groups"
		><LucideChevronLeft class="h-5 w-5" /> Back</a
	>
</div>
<Card.Root>
	<Card.Header>
		<Card.Title>Meta data</Card.Title>
	</Card.Header>

	<Card.Content>
		<UserGroupForm existingUserGroup={userGroup} callback={updateUserGroup} />
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Users</Card.Title>
		<Card.Description>Assign users to this group.</Card.Description>
	</Card.Header>

	<Card.Content>
		{#await userService.list() then users}
			<UserSelection {users} bind:selectedUserIds={userGroup.userIds} />
		{/await}
		<div class="mt-5 flex justify-end">
			<Button on:click={() => updateUserGroupUsers(userGroup.userIds)}>Save</Button>
		</div>
	</Card.Content>
</Card.Root>
