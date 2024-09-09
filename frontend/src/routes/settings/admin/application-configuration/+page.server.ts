import AppConfigService from '$lib/services/app-config-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const appConfigService = new AppConfigService(cookies.get('access_token'));
	const appConfig = await appConfigService.list(true);
	return { appConfig };
};
