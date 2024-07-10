<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { account } from '$lib/sdk';
	import { OAuthProvider } from 'appwrite';

	let email = '';
	let password = '';

	async function handleSubmit() {
		try {
			await account.createEmailPasswordSession(email, password);
			await goto('/dashboard');
		} catch (error) {
			console.log(error);
		}
	}

	async function handleGitHubLogin() {
		try {
			account.createOAuth2Session(
				OAuthProvider.Github, // provider
				'https://mail-memo.vercel.app/dashboard', // redirect here on success
				'https://mail-memo.vercel.app/login', // redirect here on failure
				['repo', 'user']
			);
		} catch (error) {
			console.log(error);
		}
	}
</script>

<div class="w-full lg:grid lg:min-h-[600px] lg:grid-cols-2 xl:min-h-[800px]">
	<div class="hidden bg-muted lg:block">
		<img
			src="/images/placeholder.svg"
			alt="placeholder"
			width="1920"
			height="1080"
			class="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
		/>
	</div>
	<div class="flex items-center justify-center py-12">
		<div class="mx-auto grid w-[350px] gap-6">
			<div class="grid gap-2 text-center">
				<h1 class="text-3xl font-bold">Register</h1>
				<p class="text-balance text-muted-foreground">Enter your email below to register</p>
			</div>
			<div class="grid gap-4">
				<div class="grid gap-2">
					<Label for="email">Email</Label>
					<Input id="email" type="email" placeholder="m@example.com" required />
				</div>
				<div class="grid gap-2">
					<div class="flex items-center">
						<Label for="password">Password</Label>
						<!-- <a href="##" class="ml-auto inline-block text-sm underline"> Forgot your password? </a> -->
					</div>
					<Input id="password" type="password" required />
				</div>
				<Button type="submit" class="w-full">Login</Button>
				<Button variant="outline" class="w-full" on:click={handleGitHubLogin}>
					Login with GitHub
				</Button>
			</div>
			<div class="mt-4 text-center text-sm">
				Already have an account?
				<a href="/login" class="underline">Sign in</a>
			</div>
		</div>
	</div>
</div>
