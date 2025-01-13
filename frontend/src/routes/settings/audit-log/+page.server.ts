import AuditLogService from '$lib/services/audit-log-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const auditLogService = new AuditLogService(cookies.get('access_token'));
	const auditLogs = await auditLogService.list({
		sort: {
			column: 'createdAt',
			direction: 'desc'
		}
	});
	return {
		auditLogs
	};
};
