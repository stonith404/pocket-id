import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, url }) => {
	return {
		token: params.token,
		redirect: url.searchParams.get('redirect') || '/settings'
	};
};
