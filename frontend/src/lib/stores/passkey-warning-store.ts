import { writable } from 'svelte/store';

const defaultPasskeyWarningVisibility = true;
const { subscribe, set: setAlertVisibility } = writable(defaultPasskeyWarningVisibility);

export const passkeyWarningVisibility = {
  subscribe,
  toggleVisibility: () => {
    const currentVisibility = localStorage.getItem('passkeyWarningVisible');
    if (currentVisibility === 'true') {
      setAlertVisibility(false);
      localStorage.setItem('passkeyWarningVisible', 'false');
    } else {
      setAlertVisibility(true);
      localStorage.setItem('passkeyWarningVisible', 'true');
    }
  },
};