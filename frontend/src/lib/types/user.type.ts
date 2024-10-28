import type { CustomClaim } from './custom-claim.type';

export type User = {
	id: string;
	username: string;
	email: string;
	firstName: string;
	lastName: string;
	isAdmin: boolean;
	customClaims: CustomClaim[];
};

export type UserCreate = Omit<User, 'id' | 'customClaims'>;
