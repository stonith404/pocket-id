<script>
	import Logo from './logo.svelte';
</script>

<div class="flex flex-col justify-center">
	<div class="mx-auto rounded-2xl bg-muted p-3">
		<Logo class="h-10 w-10" />
	</div>
	<p class="mt-5 font-playfair text-3xl font-bold sm:text-4xl">Browser unsupported</p>
	<p class="mt-3 text-muted-foreground">
		This browser doesn't support passkeys. Please use a browser that supports WebAuthn to sign in.
	</p>
</div>
