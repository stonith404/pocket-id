<script lang="ts">
	import { page } from '$app/stores';
	import { openConfirmDialog } from '$lib/components/confirm-dialog/';
	import { Badge } from '$lib/components/ui/badge/index';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import * as Pagination from '$lib/components/ui/pagination';
	import * as Table from '$lib/components/ui/table';
	import UserService from '$lib/services/user-service';
	import type { Paginated, PaginationRequest } from '$lib/types/pagination.type';
	import type { User } from '$lib/types/user.type';
	import { debounced } from '$lib/utils/debounce-util';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { LucideLink, LucidePencil, LucideTrash } from 'lucide-svelte';
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import { toast } from 'svelte-sonner';
	import OneTimeLinkModal from './one-time-link-modal.svelte';

	let { users: initialUsers }: { users: Paginated<User> } = $props();
	let users = $state<Paginated<User>>(initialUsers);
	let oneTimeLink = $state<string | null>(null);

	$effect(() => {
		users = initialUsers;
	});

	const userService = new UserService();

	let pagination = $state<PaginationRequest>({
		page: 1,
		limit: 10
	});
	let search = $state('');

	const debouncedFetchUsers = debounced(userService.list, 500);

	async function deleteUser(user: User) {
		openConfirmDialog({
			title: `Delete ${user.firstName} ${user.lastName}`,
			message: 'Are you sure you want to delete this user?',
			confirm: {
				label: 'Delete',
				destructive: true,
				action: async () => {
					try {
						await userService.remove(user.id);
						users = await userService.list(search, pagination);
					} catch (e) {
						axiosErrorToast(e);
					}
					toast.success('User deleted successfully');
				}
			}
		});
	}

	async function createOneTimeAccessToken(userId: string) {
		try {
			const token = await userService.createOneTimeAccessToken(userId);
			oneTimeLink = `${$page.url.origin}/login/${token}`;
		} catch (e) {
			axiosErrorToast(e);
		}
	}
</script>

<Input
	type="search"
	placeholder="Search users"
	bind:value={search}
	on:input={async (e) =>
		(users = await userService.list((e.target as HTMLInputElement).value, pagination))}
/>
<Table.Root>
	<Table.Header>
		<Table.Row> 
			<Table.Head class="hidden md:table-cell">First name</Table.Head>
			<Table.Head class="hidden md:table-cell">Last name</Table.Head>
			<Table.Head>Email</Table.Head>
			<Table.Head>Username</Table.Head>
			<Table.Head class="hidden lg:table-cell">Role</Table.Head>
			<Table.Head>
				<span class="sr-only">Actions</span>
			</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#if users.data.length === 0}
			<Table.Row>
				<Table.Cell colspan={6} class="text-center">No users found</Table.Cell>
			</Table.Row>
		{:else}
			{#each users.data as user}
				<Table.Row>
					<Table.Cell class="hidden md:table-cell">{user.firstName}</Table.Cell>
					<Table.Cell class="hidden md:table-cell">{user.lastName}</Table.Cell>
					<Table.Cell>{user.email}</Table.Cell>
					<Table.Cell>{user.username}</Table.Cell>
					<Table.Cell class="hidden lg:table-cell">
						<Badge variant="outline">{user.isAdmin ? 'Admin' : 'User'}</Badge>
					</Table.Cell>
					<Table.Cell>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger asChild let:builder>
								<Button aria-haspopup="true" size="icon" variant="ghost" builders={[builder]}>
									<Ellipsis class="h-4 w-4" />
									<span class="sr-only">Toggle menu</span>
								</Button>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end">
								<DropdownMenu.Item on:click={() => createOneTimeAccessToken(user.id)}
									><LucideLink class="mr-2 h-4 w-4" />One-time link</DropdownMenu.Item
								>
								<DropdownMenu.Item href="/settings/admin/users/{user.id}"
									><LucidePencil class="mr-2 h-4 w-4" /> Edit</DropdownMenu.Item
								>
								<DropdownMenu.Item
									class="text-red-500 focus:!text-red-700"
									on:click={() => deleteUser(user)}
									><LucideTrash class="mr-2 h-4 w-4" />Delete</DropdownMenu.Item
								>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</Table.Cell>
				</Table.Row>
			{/each}
		{/if}
	</Table.Body>
</Table.Root>

{#if users?.data?.length ?? 0 > 0}
	<Pagination.Root
		class="mt-5"
		count={users.pagination.totalItems}
		perPage={pagination.limit}
		onPageChange={async (p) =>
			(users = await userService.list(search, {
				page: p,
				limit: pagination.limit
			}))}
		bind:page={users.pagination.currentPage}
		let:pages
		let:currentPage
	>
		<Pagination.Content class="flex justify-end">
			<Pagination.Item>
				<Pagination.PrevButton />
			</Pagination.Item>
			{#each pages as page (page.key)}
				{#if page.type === 'ellipsis'}
					<Pagination.Item>
						<Pagination.Ellipsis />
					</Pagination.Item>
				{:else}
					<Pagination.Item>
						<Pagination.Link {page} isActive={users.pagination.currentPage === page.value}>
							{page.value}
						</Pagination.Link>
					</Pagination.Item>
				{/if}
			{/each}
			<Pagination.Item>
				<Pagination.NextButton />
			</Pagination.Item>
		</Pagination.Content>
	</Pagination.Root>
{/if}

<OneTimeLinkModal {oneTimeLink} />
