import ApplicationConfigurationService from '$lib/services/application-configuration-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const applicationConfigurationService = new ApplicationConfigurationService(
		cookies.get('access_token')
	);
	const applicationConfiguration = await applicationConfigurationService.list();
	return { applicationConfiguration };
};
