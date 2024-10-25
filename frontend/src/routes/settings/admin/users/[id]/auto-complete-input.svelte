<script lang="ts">
	import Input from '$lib/components/ui/input/input.svelte';
	import * as Popover from '$lib/components/ui/popover/index.js';

	let {
		value = $bindable(''),
		placeholder,
		suggestions
	}: { value: string; placeholder: string; suggestions: string[] } = $props();

	let suggestionResult: string[] = $state(suggestions);
	let selectedIndex = $state(-1);

	let isInputFocused = $state(false);

	function handleSuggestionClick(suggestion: (typeof suggestions)[0]) {
		value = suggestion;
		suggestionResult = [];
	}

	function handleOnInput() {
		suggestionResult = suggestions.filter((s) => s.includes(value.toLowerCase()));
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!isOpen) return;
		switch (e.key) {
			case 'ArrowDown':
				selectedIndex = Math.min(selectedIndex + 1, suggestionResult.length - 1);
				break;
			case 'ArrowUp':
				selectedIndex = Math.max(selectedIndex - 1, -1);
				break;
			case 'Enter':
				if (selectedIndex >= 0) {
					handleSuggestionClick(suggestionResult[selectedIndex]);
				}
				break;
			case 'Escape':
				isInputFocused = false;
				break;
		}
	}

	let isOpen = $derived(suggestionResult.length > 0 && isInputFocused);

	$effect(() => {
		// Reset selection when suggestions change
		if (suggestionResult) {
			selectedIndex = -1;
		}
	});
</script>

<div
	class="w-full grid"
	role="combobox"
	onkeydown={handleKeydown}
	aria-controls="suggestion-list"
	aria-expanded={isOpen}
	tabindex="-1"
>
	<Input
		{placeholder}
		bind:value
		oninput={handleOnInput}
		onfocus={() => (isInputFocused = true)}
		onblur={() => (isInputFocused = false)}
	/>
	<Popover.Root
		open={isOpen}
		disableFocusTrap
		openFocus={() => {}}
		closeOnOutsideClick={false}
		closeOnEscape={false}
	>
		<Popover.Trigger tabindex={-1} class="w-full h-0" aria-hidden />
		<Popover.Content class="p-0" sideOffset={5} sameWidth>
			{#each suggestionResult as suggestion, index}
				<div
					role="button"
					tabindex="0"
					onmousedown={() => handleSuggestionClick(suggestion)}
					onkeydown={(e) => {
						if (e.key === 'Enter') handleSuggestionClick(suggestion);
					}}
					class="hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50 {selectedIndex ===
					index
						? 'bg-accent text-accent-foreground'
						: ''}"
				>
					{suggestion}
				</div>
			{/each}
		</Popover.Content>
	</Popover.Root>
</div>
