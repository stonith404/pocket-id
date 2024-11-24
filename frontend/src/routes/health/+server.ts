import AppConfigService from '$lib/services/app-config-service';
import type { RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async () => {
	const appConfigService = new AppConfigService();
	let backendOk = true;
	await appConfigService.list().catch(() => (backendOk = false));

	return new Response(
		JSON.stringify({
			status: backendOk ? 'HEALTHY' : 'UNHEALTHY'
		}),
		{
			status: backendOk ? 200 : 500,
			headers: {
				'content-type': 'application/json'
			}
		}
	);
};
