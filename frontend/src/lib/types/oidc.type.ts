export type OidcClient = {
	id: string;
	name: string;
	logoURL: string;
	callbackURL: string;
	hasLogo: boolean;
};

export type OidcClientCreate = Omit<OidcClient, 'id' | 'logoURL' | 'hasLogo'>;

export type OidcClientCreateWithLogo = OidcClientCreate & {
	logo: File | null;
};
