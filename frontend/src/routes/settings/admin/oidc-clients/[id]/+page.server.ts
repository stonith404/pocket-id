import { ACCESS_TOKEN_COOKIE_NAME } from '$lib/constants';
import OidcService from '$lib/services/oidc-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const oidcService = new OidcService(cookies.get(ACCESS_TOKEN_COOKIE_NAME));
	return await oidcService.getClient(params.id);
};
