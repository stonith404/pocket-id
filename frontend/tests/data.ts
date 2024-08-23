export const users = {
	tim: {
		id: 'f4b89dc2-62fb-46bf-9f5f-c34f4eafe93e',
		firstname: 'Tim',
		lastname: 'Cook',
		email: 'tim.cook@test.com',
		username: 'tim'
	},
	craig: {
		id: '1cd19686-f9a6-43f4-a41f-14a0bf5b4036',
		firstname: 'Craig',
		lastname: 'Federighi',
		email: 'craig.federighi@test.com',
		username: 'craig'
	},
	steve: {
		firstname: 'Steve',
		lastname: 'Jobs',
		email: 'steve.jobs@test.com',
		username: 'steve'
	}
};

export const oidcClients = {
	nextcloud: {
		id: '3654a746-35d4-4321-ac61-0bdcff2b4055',
		name: 'Nextcloud',
		callbackUrl: 'http://nextcloud/auth/callback'
	},
	immich: {
		id: '606c7782-f2b1-49e5-8ea9-26eb1b06d018',
		name: 'Immich',
		callbackUrl: 'http://immich/auth/callback'
	},
	pingvinShare: {
		name: 'Pingvin Share',
		callbackUrl: 'http://pingvin.share/auth/callback',
		secondCallbackUrl: 'http://pingvin.share/auth/callback2'
	}
};
