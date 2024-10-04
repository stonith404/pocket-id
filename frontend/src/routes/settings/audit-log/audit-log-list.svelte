<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Pagination from '$lib/components/ui/pagination';
	import * as Table from '$lib/components/ui/table';
	import AuditLogService from '$lib/services/audit-log-service';
	import type { AuditLog } from '$lib/types/audit-log.type';
	import type { Paginated, PaginationRequest } from '$lib/types/pagination.type';

	let { auditLogs: initialAuditLog }: { auditLogs: Paginated<AuditLog> } = $props();
	let auditLogs = $state<Paginated<AuditLog>>(initialAuditLog);

	const auditLogService = new AuditLogService();

	let pagination = $state<PaginationRequest>({
		page: 1,
		limit: 15
	});

	function toFriendlyEventString(event: string) {
		const words = event.split('_');
		const capitalizedWords = words.map((word) => {
			return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
		});
		return capitalizedWords.join(' ');
	}
</script>

<Table.Root>
	<Table.Header class="whitespace-nowrap">
		<Table.Row>
			<Table.Head>Time</Table.Head>
			<Table.Head>Event</Table.Head>
			<Table.Head>Approximate Location</Table.Head>
			<Table.Head>IP Address</Table.Head>
			<Table.Head>Device</Table.Head>
			<Table.Head>Client</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body class="whitespace-nowrap">
		{#if auditLogs.data.length === 0}
			<Table.Row>
				<Table.Cell colspan={6} class="text-center">No logs found</Table.Cell>
			</Table.Row>
		{:else}
			{#each auditLogs.data as auditLog}
				<Table.Row>
					<Table.Cell>{new Date(auditLog.createdAt).toLocaleString()}</Table.Cell>
					<Table.Cell>
						<Badge variant="outline">{toFriendlyEventString(auditLog.event)}</Badge>
					</Table.Cell>
					<Table.Cell>{auditLog.city && auditLog.country ? `${auditLog.city}, ${auditLog.country}` : 'Unknown'}</Table.Cell>
					<Table.Cell>{auditLog.ipAddress}</Table.Cell>
					<Table.Cell>{auditLog.device}</Table.Cell>
					<Table.Cell>{auditLog.data.clientName}</Table.Cell>
				</Table.Row>
			{/each}
		{/if}
	</Table.Body>
</Table.Root>

{#if auditLogs?.data?.length ?? 0 > 0}
	<Pagination.Root
		class="mt-5"
		count={auditLogs.pagination.totalItems}
		perPage={pagination.limit}
		onPageChange={async (p) =>
			(auditLogs = await auditLogService.list({
				page: p,
				limit: pagination.limit
			}))}
		bind:page={auditLogs.pagination.currentPage}
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
						<Pagination.Link {page} isActive={auditLogs.pagination.currentPage === page.value}>
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
