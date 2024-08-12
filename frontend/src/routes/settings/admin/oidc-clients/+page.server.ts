import OIDCService from '$lib/services/oidc-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const oidcService = new OIDCService(cookies.get('access_token'));
	const clients = await oidcService.listClients();
	return clients;
};
