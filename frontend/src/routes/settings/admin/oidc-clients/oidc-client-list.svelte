<script lang="ts">
	import { openConfirmDialog } from '$lib/components/confirm-dialog/';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Pagination from '$lib/components/ui/pagination';
	import * as Table from '$lib/components/ui/table';
	import OIDCService from '$lib/services/oidc-service';
	import type { OidcClient } from '$lib/types/oidc.type';
	import type { Paginated, PaginationRequest } from '$lib/types/pagination.type';
	import { debounced } from '$lib/utils/debounce-util';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { LucidePencil, LucideTrash } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import OneTimeLinkModal from './client-secret.svelte';

	let { clients: initialClients }: { clients: Paginated<OidcClient> } = $props();
	let clients = $state<Paginated<OidcClient>>(initialClients);
	let oneTimeLink = $state<string | null>(null);

	$effect(() => {
		clients = initialClients;
	});

	const oidcService = new OIDCService();

	let pagination = $state<PaginationRequest>({
		page: 1,
		limit: 10
	});
	let search = $state('');

	const debouncedSearch = debounced(async (searchValue: string) => {
		clients = await oidcService.listClients(searchValue, pagination);
	}, 400);

	async function deleteClient(client: OidcClient) {
		openConfirmDialog({
			title: `Delete ${client.name}`,
			message: 'Are you sure you want to delete this OIDC client?',
			confirm: {
				label: 'Delete',
				destructive: true,
				action: async () => {
					try {
						await oidcService.removeClient(client.id);
						clients = await oidcService.listClients(search, pagination);
						toast.success('OIDC client deleted successfully');
					} catch (e) {
						axiosErrorToast(e);
					}
				}
			}
		});
	}
</script>

<Input
	type="search"
	placeholder="Search clients"
	bind:value={search}
	on:input={(e) => debouncedSearch((e.target as HTMLInputElement).value)}
/>
<Table.Root>
	<Table.Header class="sr-only">
		<Table.Row>
			<Table.Head>Logo</Table.Head>
			<Table.Head>Name</Table.Head>
			<Table.Head>Actions</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#if clients.data.length === 0}
			<Table.Row>
				<Table.Cell colspan={6} class="text-center">No OIDC clients found</Table.Cell>
			</Table.Row>
		{:else}
			{#each clients.data as client}
				<Table.Row>
					<Table.Cell class="w-8 font-medium">
						{#if client.hasLogo}
							<div class="h-8 w-8">
								<img
									class="m-auto max-h-full max-w-full object-contain"
									src="/api/oidc/clients/{client.id}/logo"
									alt="{client.name} logo"
								/>
							</div>
						{/if}
					</Table.Cell>
					<Table.Cell class="font-medium">{client.name}</Table.Cell>
					<Table.Cell class="flex justify-end gap-1">
						<Button
							href="/settings/admin/oidc-clients/{client.id}"
							size="sm"
							variant="outline"
							aria-label="Edit"><LucidePencil class="h-3 w-3 " /></Button
						>
						<Button
							on:click={() => deleteClient(client)}
							size="sm"
							variant="outline"
							aria-label="Delete"><LucideTrash class="h-3 w-3 text-red-500" /></Button
						>
					</Table.Cell>
				</Table.Row>
			{/each}
		{/if}
	</Table.Body>
</Table.Root>

{#if clients?.data?.length ?? 0 > 0}
	<Pagination.Root
		class="mt-5"
		count={clients.pagination.totalItems}
		perPage={pagination.limit}
		onPageChange={async (p) =>
			(clients = await oidcService.listClients(search, {
				page: p,
				limit: pagination.limit
			}))}
		bind:page={clients.pagination.currentPage}
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
						<Pagination.Link {page} isActive={clients.pagination.currentPage === page.value}>
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
