import { ACCESS_TOKEN_COOKIE_NAME } from '$lib/constants';
import AppConfigService from '$lib/services/app-config-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const appConfigService = new AppConfigService(cookies.get(ACCESS_TOKEN_COOKIE_NAME));
	const appConfig = await appConfigService.list(true);
	return { appConfig };
};
