import { writable } from 'svelte/store';

const clientSecretStore = writable<string | null>(null);

const set = (user: string) => {
	clientSecretStore.set(user);
};

const clear = () => {
	clientSecretStore.set(null);
};

export default {
	subscribe: clientSecretStore.subscribe,
	set,
	clear
};
