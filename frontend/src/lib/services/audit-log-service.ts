import type { AuditLog } from '$lib/types/audit-log.type';
import type { Paginated, PaginationRequest } from '$lib/types/pagination.type';
import APIService from './api-service';

class AuditLogService extends APIService {
	async list(pagination?: PaginationRequest) {
		const page = pagination?.page || 1;
		const limit = pagination?.limit || 10;

		const res = await this.api.get('/audit-logs', {
			params: {
				page,
				limit
			}
		});
		return res.data as Paginated<AuditLog>;
	}
}

export default AuditLogService;
