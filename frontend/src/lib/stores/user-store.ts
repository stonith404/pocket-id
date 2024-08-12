import type { User } from '$lib/types/user.type';
import { writable } from 'svelte/store';

const userStore = writable<User | null>(null);

const setUser = (user: User) => {
	userStore.set(user);
};

const clearUser = () => {
	userStore.set(null);
};

export default {
	subscribe: userStore.subscribe,
	setUser,
	clearUser
};
