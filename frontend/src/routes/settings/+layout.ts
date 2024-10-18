import AppConfigService from '$lib/services/app-config-service';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
	const appConfigService = new AppConfigService();

	const versionInformation = await appConfigService.getVersionInformation();
	return {
		versionInformation
	};
};
