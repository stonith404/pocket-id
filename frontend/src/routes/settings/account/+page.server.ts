import UserService from '$lib/services/user-service';
import WebAuthnService from '$lib/services/webauthn-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const webauthnService = new WebAuthnService(cookies.get('access_token'));
	const userService = new UserService(cookies.get('access_token'));
	const account = await userService.getCurrent();
	const passkeys = await webauthnService.listCredentials();
	return {
		account,
		passkeys
	};
};
