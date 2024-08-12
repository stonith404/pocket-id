import { WebAuthnError } from '@simplewebauthn/browser';
import { AxiosError } from 'axios';
import { toast } from 'svelte-sonner';

export function axiosErrorToast(e: unknown, message: string = 'An unknown error occurred') {
	if (e instanceof AxiosError) {
		message = e.response?.data.error || message;
	}
	toast.error(message);
}

export function getWebauthnErrorMessage(e: unknown) {
	const errors = {
		ERROR_CEREMONY_ABORTED: 'The authentication process was aborted',
		ERROR_AUTHENTICATOR_GENERAL_ERROR: 'An error occurred with the authenticator',
		ERROR_AUTHENTICATOR_MISSING_DISCOVERABLE_CREDENTIAL_SUPPORT:
			'The authenticator does not support discoverable credentials',
		ERROR_AUTHENTICATOR_MISSING_RESIDENT_KEY_SUPPORT:
			'The authenticator does not support resident keys',
		ERROR_AUTHENTICATOR_PREVIOUSLY_REGISTERED: 'This passkey was previously registered',
		ERROR_AUTHENTICATOR_NO_SUPPORTED_PUBKEYCREDPARAMS_ALG:
			'The authenticator does not support any of the requested algorithms'
	};

	let message = 'An unknown error occurred';
	if (e instanceof WebAuthnError && e.code in errors) {
		message = errors[e.code as keyof typeof errors];
	} else if (e instanceof WebAuthnError && e?.message.includes('timed out')) {
		message = 'The authenticator timed out';
	} else if (e instanceof AxiosError && e.response?.data.error) {
		message = e.response?.data.error;
	} else {
		console.error(e);
	}
	return message;
}
