export type OidcClient = {
	id: string;
	name: string;
	logoURL: string;
	callbackURLs: [string, ...string[]];
	hasLogo: boolean;
	isPublic: boolean;
};

export type OidcClientCreate = Omit<OidcClient, 'id' | 'logoURL' | 'hasLogo'>;

export type OidcClientCreateWithLogo = OidcClientCreate & {
	logo: File | null;
};

export type AuthorizeResponse = {
	code: string;
	callbackURL: string;
};
