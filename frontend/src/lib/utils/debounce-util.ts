export function debounced<T extends (...args: any[]) => void>(func: T, delay: number) {
	let debounceTimeout: ReturnType<typeof setTimeout>;

	return (...args: Parameters<T>) => {
		if (debounceTimeout !== undefined) {
			clearTimeout(debounceTimeout);
		}

		debounceTimeout = setTimeout(() => {
			func(...args);
		}, delay);
	};
}
