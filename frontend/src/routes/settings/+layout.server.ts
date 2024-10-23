import AppConfigService from '$lib/services/app-config-service';
import type { AppVersionInformation } from '$lib/types/application-configuration';
import type { LayoutServerLoad } from './$types';

let versionInformation: AppVersionInformation;
let versionInformationLastUpdated: number;

export const load: LayoutServerLoad = async () => {
	const appConfigService = new AppConfigService();

	// Cache the version information for 3 hours
	const cacheExpired =
		versionInformationLastUpdated &&
		Date.now() - versionInformationLastUpdated > 1000 * 60 * 60 * 3;

	if (!versionInformation || cacheExpired) {
		versionInformation = await appConfigService.getVersionInformation();
		versionInformationLastUpdated = Date.now();
	}

	return {
		versionInformation
	};
};
