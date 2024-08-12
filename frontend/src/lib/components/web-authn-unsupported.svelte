<script>
	import Logo from './logo.svelte';
</script>

<div class="flex flex-col justify-center">
	<div class="bg-muted mx-auto rounded-2xl p-3">
		<Logo class="h-10 w-10" />
	</div>
	<p class="font-playfair mt-5 text-3xl sm:text-4xl font-bold">Browser unsupported</p>
	<p class="text-muted-foreground mt-3">
		This browser doesn't support passkeys. Please use a browser that supports WebAuthn to sign in.
	</p>
</div>
