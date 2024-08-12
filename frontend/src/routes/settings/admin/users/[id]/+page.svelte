<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import UserService from '$lib/services/user-service';
	import type { UserCreate } from '$lib/types/user.type';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { LucideChevronLeft } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import UserForm from '../user-form.svelte';

	let { data } = $props();
	let user = $state(data);

	const userService = new UserService();

	async function updateUser(updatedUser: UserCreate) {
		let success = true;
		await userService
			.update(user.id, updatedUser)
			.then(() => toast.success('User updated successfully'))
			.catch((e) => {
				axiosErrorToast(e);
				success = false;
			});

		return success;
	}
</script>

<svelte:head>
	<title>User Details {user.firstName} {user.lastName}</title>
</svelte:head>

<div>
	<a class="text-muted-foreground flex text-sm" href="/settings/admin/users"
		><LucideChevronLeft class="h-5 w-5" /> Back</a
	>
</div>
<Card.Root>
	<Card.Header>
		<Card.Title>{user.firstName} {user.lastName}</Card.Title>
	</Card.Header>

	<Card.Content>
		<UserForm existingUser={user} callback={updateUser} />
	</Card.Content>
</Card.Root>
