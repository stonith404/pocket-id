import { ACCESS_TOKEN_COOKIE_NAME } from '$lib/constants';
import AuditLogService from '$lib/services/audit-log-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const auditLogService = new AuditLogService(cookies.get(ACCESS_TOKEN_COOKIE_NAME));
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
