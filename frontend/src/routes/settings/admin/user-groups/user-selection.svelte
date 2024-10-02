<script lang="ts">
	import AdvancedTable from '$lib/components/advanced-table.svelte';
	import * as Table from '$lib/components/ui/table';
	import UserService from '$lib/services/user-service';
	import type { Paginated } from '$lib/types/pagination.type';
	import type { User } from '$lib/types/user.type';

	let {
		users: initialUsers,
		selectedUserIds = $bindable()
	}: { users: Paginated<User>; selectedUserIds: string[] } = $props();

	const userService = new UserService();

	let users = $state(initialUsers);

	function fetchItems(search: string, page: number, limit: number) {
		return userService.list(search, { page, limit });
	}
</script>

<AdvancedTable
	items={users}
	{fetchItems}
	columns={['Name', 'Email']}
	bind:selectedIds={selectedUserIds}
>
	{#snippet rows({ item })}
		<Table.Cell>{item.firstName} {item.lastName}</Table.Cell>
		<Table.Cell>{item.email}</Table.Cell>
	{/snippet}
</AdvancedTable>
