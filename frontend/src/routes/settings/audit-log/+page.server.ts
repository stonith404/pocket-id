import AuditLogService from '$lib/services/audit-log-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const auditLogService = new AuditLogService(cookies.get('access_token'));
	const auditLogs = await auditLogService.list({
		limit: 15,
		page: 1,
	});
	return {
		auditLogs
	};
};
