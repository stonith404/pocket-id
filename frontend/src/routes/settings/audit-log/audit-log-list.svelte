<script lang="ts">
	import AdvancedTable from '$lib/components/advanced-table.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import AuditLogService from '$lib/services/audit-log-service';
	import type { AuditLog } from '$lib/types/audit-log.type';
	import type { Paginated } from '$lib/types/pagination.type';

	let { auditLogs: initialAuditLog }: { auditLogs: Paginated<AuditLog> } = $props();
	let auditLogs = $state<Paginated<AuditLog>>(initialAuditLog);

	const auditLogService = new AuditLogService();

	async function fetchItems(search: string, page: number, limit: number) {
		return await auditLogService.list({
			page,
			limit
		});
	}

	function toFriendlyEventString(event: string) {
		const words = event.split('_');
		const capitalizedWords = words.map((word) => {
			return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
		});
		return capitalizedWords.join(' ');
	}
</script>

<AdvancedTable
	items={auditLogs}
	{fetchItems}
	columns={['Time', 'Event', 'Approximate Location', 'IP Address', 'Device', 'Client']}
	withoutSearch
>
	{#snippet rows({ item })}
		<Table.Cell>{new Date(item.createdAt).toLocaleString()}</Table.Cell>
		<Table.Cell>
			<Badge variant="outline">{toFriendlyEventString(item.event)}</Badge>
		</Table.Cell>
		<Table.Cell
			>{item.city && item.country ? `${item.city}, ${item.country}` : 'Unknown'}</Table.Cell
		>
		<Table.Cell>{item.ipAddress}</Table.Cell>
		<Table.Cell>{item.device}</Table.Cell>
		<Table.Cell>{item.data.clientName}</Table.Cell>
	{/snippet}
</AdvancedTable>
