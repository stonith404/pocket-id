<script lang="ts" generics="T extends {id:string}">
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Pagination from '$lib/components/ui/pagination';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table/index.js';
	import Empty from '$lib/icons/empty.svelte';
	import type { Paginated } from '$lib/types/pagination.type';
	import { debounced } from '$lib/utils/debounce-util';
	import type { Snippet } from 'svelte';

	let {
		items,
		selectedIds = $bindable(),
		withoutSearch = false,
		fetchItems,
		columns,
		rows
	}: {
		items: Paginated<T>;
		selectedIds?: string[];
		withoutSearch?: boolean;
		fetchItems: (search: string, page: number, limit: number) => Promise<Paginated<T>>;
		columns: (string | { label: string; hidden?: boolean })[];
		rows: Snippet<[{ item: T }]>;
	} = $props();

	let availablePageSizes: number[] = [10, 20, 50, 100];

	let allChecked = $derived.by(() => {
		if (!selectedIds || items.data.length === 0) return false;
		for (const item of items.data) {
			if (!selectedIds.includes(item.id)) {
				return false;
			}
		}
		return true;
	});

	const onSearch = debounced(async (searchValue: string) => {
		items = await fetchItems(searchValue, 1, items.pagination.itemsPerPage);
	}, 300);

	async function onAllCheck(checked: boolean) {
		if (checked) {
			selectedIds = items.data.map((item) => item.id);
		} else {
			selectedIds = [];
		}
	}

	async function onCheck(checked: boolean, id: string) {
		if (!selectedIds) return;
		if (checked) {
			selectedIds = [...selectedIds, id];
		} else {
			selectedIds = selectedIds.filter((selectedId) => selectedId !== id);
		}
	}

	async function onPageChange(page: number) {
		items = await fetchItems('', page, items.pagination.itemsPerPage);
	}

	async function onPageSizeChange(size: number) {
		items = await fetchItems('', 1, size);
	}
</script>

{#if items.data.length === 0}
	<div class="my-5 flex flex-col items-center">
		<Empty class="text-muted-foreground h-20" />
		<p class="text-muted-foreground mt-3 text-sm">No items found</p>
	</div>
{:else}
	<div class="w-full">
		{#if !withoutSearch}
			<Input
				class="mb-4 max-w-sm"
				placeholder={'Search...'}
				type="text"
				oninput={(e) => onSearch((e.target as HTMLInputElement).value)}
			/>
		{/if}

		<Table.Root>
			<Table.Header>
				<Table.Row>
					{#if selectedIds}
						<Table.Head>
							<Checkbox checked={allChecked} onCheckedChange={(c) => onAllCheck(c as boolean)} />
						</Table.Head>
					{/if}
					{#each columns as column}
						{#if typeof column === 'string'}
							<Table.Head>{column}</Table.Head>
						{:else}
							<Table.Head class={column.hidden ? 'sr-only' : ''}>{column.label}</Table.Head>
						{/if}
					{/each}
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each items.data as item}
					<Table.Row class={selectedIds?.includes(item.id) ? 'bg-muted/20' : ''}>
						{#if selectedIds}
							<Table.Cell>
								<Checkbox
									checked={selectedIds.includes(item.id)}
									onCheckedChange={(c) => onCheck(c as boolean, item.id)}
								/>
							</Table.Cell>
						{/if}
						{@render rows({ item })}
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>

		<div class="mt-5 flex items-center justify-between space-x-2">
			<div class="flex items-center space-x-2">
				<p class="text-sm font-medium">Items per page</p>
				<Select.Root
					selected={{
						label: items.pagination.itemsPerPage.toString(),
						value: items.pagination.itemsPerPage
					}}
					onSelectedChange={(v) => onPageSizeChange(v?.value as number)}
				>
					<Select.Trigger class="h-9 w-[80px]">
						<Select.Value>{items.pagination.itemsPerPage}</Select.Value>
					</Select.Trigger>
					<Select.Content>
						{#each availablePageSizes as size}
							<Select.Item value={size}>{size}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
			<Pagination.Root
				class="mx-0 w-auto"
				count={items.pagination.totalItems}
				perPage={items.pagination.itemsPerPage}
				{onPageChange}
				page={items.pagination.currentPage}
				let:pages
			>
				<Pagination.Content class="flex justify-end">
					<Pagination.Item>
						<Pagination.PrevButton />
					</Pagination.Item>
					{#each pages as page (page.key)}
						{#if page.type !== 'ellipsis'}
							<Pagination.Item>
								<Pagination.Link {page} isActive={items.pagination.currentPage === page.value}>
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
		</div>
	</div>
{/if}
