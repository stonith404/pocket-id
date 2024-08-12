export type User = {
	id: string;
	username: string;
	email: string;
	firstName: string;
	lastName: string;
	isAdmin: boolean;
};

export type UserCreate = Omit<User, 'id'>;
