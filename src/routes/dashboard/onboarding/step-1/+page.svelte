<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { onMount } from 'svelte';
	import { step } from '../store';
	import { databases, functions } from '$lib/sdk';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';
	import { ExecutionMethod, ID } from 'appwrite';
	import { CalColors, colorToHex } from '$lib/calendars';
	import { toast } from 'svelte-sonner';

	let selectedProvider: 'google' | 'apple' | 'outlook' | 'url' | null = null;
	let name = '';
	let url = '';
	let color = '';
	let form: HTMLFormElement;

	onMount(() => {
		step.set(1);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
			return;
		}
		if (selectedProvider === 'url') {
			try {
				const execution = await functions.createExecution(
					'api',
					JSON.stringify({
						name,
						url,
						color
					}),
					false,
					'/v1/scheduler/intervals',
					ExecutionMethod.PATCH
				);
				const isOk = execution.responseStatusCode;

				if (!isOk) {
					toast(
						execution.responseBody ? execution.responseBody : 'Unexpected error. Please try again.'
					);
					return;
				} else {
					await goto('/dashboard/onboarding/step-2');
				}
			} catch (error) {
				console.log(error);
				const message = (error as Error)?.message;
				toast(message);
			}
		}
	}

	$: if (!selectedProvider) {
		name = '';
		url = '';
		color = '';
	}
</script>

<svelte:head>
	<title>Onboarding (step 1) - MailMemo</title>
</svelte:head>

<div>
	<h1 class="font-header mt-6 max-w-[80%] text-3xl tracking-tight">
		Connect a calendar to get started
	</h1>
	<p class="mt-4 text-muted-foreground">
		Link your preferred calendar to seamlessly manage all your important events and reminders.
	</p>
	{#if !selectedProvider}
		<div class="mt-20 grid grid-cols-2 grid-rows-2 items-center justify-center gap-4">
			<Card.Root class="opacity-50">
				<Card.Header class="p-4">
					<div>
						<Badge variant="secondary">Soon</Badge>
					</div>
				</Card.Header>
				<Card.Content></Card.Content>
				<Card.Footer class="mt-auto flex flex-col items-start justify-start gap-2 p-4">
					<img src="/icons/google.svg" alt="google" />

					<p>Google calendar</p>
				</Card.Footer>
			</Card.Root>
			<Card.Root class="opacity-50">
				<Card.Header class="p-4">
					<div>
						<Badge variant="secondary">Soon</Badge>
					</div>
				</Card.Header>
				<Card.Content></Card.Content>
				<Card.Footer class="mt-auto flex flex-col items-start justify-start gap-2 p-4">
					<img src="/icons/outlook.svg" alt="outlook" />

					<p>Outlook</p>
				</Card.Footer>
			</Card.Root>

			<Card.Root class="opacity-50">
				<Card.Header class="p-4">
					<div>
						<Badge variant="secondary">Soon</Badge>
					</div>
				</Card.Header>
				<Card.Content></Card.Content>
				<Card.Footer class="mt-auto flex flex-col items-start justify-start gap-2 p-4">
					<img src="/icons/apple.svg" alt="apple" />

					<p>Apple calendar</p>
				</Card.Footer>
			</Card.Root>
			<button on:click={() => (selectedProvider = 'url')}>
				<Card.Root>
					<Card.Header class="p-7"></Card.Header>
					<Card.Content></Card.Content>
					<Card.Footer class="mt-auto flex flex-col items-start justify-start gap-2 p-4">
						<img src="/icons/url.svg" alt="url" />

						<p>Calendar URL</p>
					</Card.Footer>
				</Card.Root>
			</button>
		</div>
	{:else if selectedProvider === 'url'}
		<form class="mt-20 grid gap-10" on:submit|preventDefault bind:this={form}>
			<div class="grid gap-2">
				<Label for="url">Calendar URL</Label>
				<Input id="url" type="url" placeholder="Enter URL" required bind:value={url} />
				<p class="text-sm text-muted-foreground">
					You can find iCal address in sharing options in your calendar settings. <a
						href="https://abalone-swing-bf0.notion.site/Instructions-to-get-iCal-link-07d2282244cb4760bdef8f8f0b1833f5"
						target="_blank"
						rel="noopener noreferrer"
						class="underline">Learn more</a
					>
				</p>
			</div>
			<div class="grid gap-2">
				<Label for="url">Calendar name</Label>
				<Input id="name" type="text" required placeholder="Enter name" bind:value={name} />
			</div>
			<div class="grid gap-2">
				<Label for="url">Calendar color</Label>
				<div class="mt-2 flex gap-1">
					{#each CalColors as c}
						<button
							type="button"
							class={`border-separate rounded-full border text-primary ring-offset-background focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:border-primary `}
							disabled={color === c}
							on:click={() => (color = c)}
						>
							<div
								class="m-1 h-10 w-10 rounded-full"
								style={`background-color: ${colorToHex(c)}`}
							/>
						</button>
					{/each}
				</div>
			</div>
		</form>
	{/if}
</div>

<div class="mt-auto flex justify-between">
	{#if selectedProvider}
		<Button variant="outline" on:click={() => (selectedProvider = null)}>Back</Button>
		<Button on:click={handleNext}>Next</Button>
	{/if}
</div>
