import OidcService from '$lib/services/oidc-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const oidcService = new OidcService(cookies.get('access_token'));
	return await oidcService.getClient(params.id);
};
