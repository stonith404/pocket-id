import type { CDPSession, Page } from '@playwright/test';

// The existing passkeys are already stored in the database
const passkeys = {
	existing1: {
		credentialId: 'test-credential-1',
		privateKey:
			'MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg3rNKkGApsEA1TpGiphKh6axTq3Vh6wBghLLea/YkIp+hRANCAATBw6jkpXXr0pHrtAQetxiR5cTcILG/YGDCdKrhVhNDHIu12YrF6B7Frwl3AUqEpdrYEwj3Fo3XkGgvrBIJEUmG'
	},
	existing2: {
		credentialId: 'test-credential-2',
		privateKey:
			'MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg3rNKkGApsEA1TpGiphKh6axTq3Vh6wBghLLea/YkIp+hRANCAATBw6jkpXXr0pHrtAQetxiR5cTcILG/YGDCdKrhVhNDHIu12YrF6B7Frwl3AUqEpdrYEwj3Fo3XkGgvrBIJEUmG'
	},
	new: {
		credentialId: 'new-test-credential',
		privateKey:
			'MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgFl2lIlRyc2G7O9D8WWrw2N8D7NTlhgWcKFY7jYxrfcmhRANCAASmvbCFrXshUvW7avTIysV9UymbhmUwGb7AonUMQPgqK2Jur7PWp9V0AIe5YMuXYH1oxsqY5CoAbdY2YsPmhYoX'
	}
};

async function init(page: Page) {
	const client = await page.context().newCDPSession(page);
	await client.send('WebAuthn.enable');
	const authenticatorId = await addVirtualAuthenticator(client);

	return {
		addPasskey: async (passkey?: keyof typeof passkeys) => {
			await addPasskey(authenticatorId, client, passkey);
		}
	};
}

async function addVirtualAuthenticator(client: CDPSession): Promise<string> {
	const result = await client.send('WebAuthn.addVirtualAuthenticator', {
		// config authenticator
		options: {
			protocol: 'ctap2',
			transport: 'internal',
			hasResidentKey: true,
			hasUserVerification: true,
			isUserVerified: true
		}
	});
	return result.authenticatorId;
}

async function addPasskey(
	authenticatorId: string,
	client: CDPSession,
	passkeyName?: keyof typeof passkeys
): Promise<void> {
	const passkey = passkeys[passkeyName ?? 'existing1'];
	await client.send('WebAuthn.addCredential', {
		authenticatorId,
		credential: {
			credentialId: btoa(passkey.credentialId),
			isResidentCredential: true,
			rpId: 'localhost',
			privateKey: passkey.privateKey,
			userHandle: btoa('f4b89dc2-62fb-46bf-9f5f-c34f4eafe93e'),
			signCount: Math.round((new Date().getTime() - 1704444610871) / 1000 / 2)
			// signCount: 2,
		}
	});
}

export default { init };
