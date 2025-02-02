export type AuditLog = {
	id: string;
	event: string;
	ipAddress: string;
	country?: string;
	city?: string;
	isp?: string;
	asNumber?: number;
	device: string;
	createdAt: string;
	data: any;
};
