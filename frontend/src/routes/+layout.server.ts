import AppConfigService from '$lib/services/app-config-service';
import UserService from '$lib/services/user-service';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	const userService = new UserService(cookies.get('access_token'));
	const appConfigService = new AppConfigService(cookies.get('access_token'));

	const user = await userService
		.getCurrent()
		.then((user) => user)
		.catch(() => null);

	const appConfig = await appConfigService
		.list()
		.then((config) => config)
		.catch((e) => {
			console.error(
				`Failed to get application configuration: ${e.response?.data.error || e.message}`
			);
			return null;
		});
	return {
		user,
		appConfig
	};
};
