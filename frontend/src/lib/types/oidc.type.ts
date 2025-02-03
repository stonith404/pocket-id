import type { UserGroup } from './user-group.type';

export type OidcClient = {
	id: string;
	name: string;
	logoURL: string;
	callbackURLs: [string, ...string[]];
	hasLogo: boolean;
	isPublic: boolean;
	pkceEnabled: boolean;
};

export type OidcClientWithAllowedUserGroups = OidcClient & {
	allowedUserGroups: UserGroup[];
};

export type OidcClientCreate = Omit<OidcClient, 'id' | 'logoURL' | 'hasLogo'>;

export type OidcClientCreateWithLogo = OidcClientCreate & {
	logo: File | null | undefined;
};

export type AuthorizeResponse = {
	code: string;
	callbackURL: string;
};
