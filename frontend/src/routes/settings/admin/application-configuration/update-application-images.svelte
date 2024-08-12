<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import ApplicationImage from './application-image.svelte';

	let {
		callback
	}: {
		callback: (logo: File | null, backgroundImage: File | null, favicon: File | null) => void;
	} = $props();

	let logo = $state<File | null>(null);
	let backgroundImage = $state<File | null>(null);
	let favicon = $state<File | null>(null);
</script>

<div class="application-images-grid">
	<ApplicationImage
		id="favicon"
		imageClass="h-14 w-14 p-2"
		label="Favicon"
		bind:image={favicon}
		imageURL="/api/application-configuration/favicon"
		accept="image/x-icon"
	/>
	<ApplicationImage
		id="logo"
		imageClass="h-32 w-32"
		label="Logo"
		bind:image={logo}
		imageURL="/api/application-configuration/logo"
	/>
	<ApplicationImage
		id="background-image"
		class="basis-full lg:basis-auto"
		imageClass="h-[350px] max-w-[500px]"
		label="Background Image"
		bind:image={backgroundImage}
		imageURL="/api/application-configuration/background-image"
	/>
</div>
<div class="flex justify-end">
	<Button class="mt-5" onclick={() => callback(logo, backgroundImage, favicon)}>Save</Button>
</div>
