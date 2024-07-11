<script lang="ts">
	import { user } from '$lib/stores';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
	import * as Card from '$lib/components/ui/card';
	import { toLocaleTimeISO } from '$lib/utils';
	import EmptyCard from './emptyCard.svelte';
	import type { CalEvent } from '$lib/calendars';
	import Button from '$lib/components/ui/button/button.svelte';

	export let data;

	function evaluateCornerRadius(index: number, lenght: number) {
		if (lenght === 1) {
			return 'rounded-xl';
		} else {
			//Last element
			if (index === lenght - 1) {
				return 'rounded-b-xl rounded-t-none';
			}
			//First element
			else if (index === 0) {
				return 'rounded-t-xl rounded-b-none';
			} else {
				return 'rounded-none';
			}
		}
	}
	$: todayEvents = data?.events?.documents?.filter((event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		return eventDate.getDate() === today.getDate();
	});

	$: futureEvents = data?.events?.documents?.filter((event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		return eventDate.getDate() > today.getDate();
	});
	//Group events that are in the same day. Use the number of days from today as the key
	$: futureEventsGrouped = futureEvents?.reduce((acc: Record<string, CalEvent[]>, event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		const days = eventDate.getDate() - today.getDate();
		if (!acc[days]) {
			acc[days] = [];
		}
		acc[days].push(event as CalEvent);
		return acc;
	}, {});

	$: console.log(futureEvents, futureEventsGrouped, Object.entries(futureEventsGrouped));
</script>

<svelte:head>
	<title>Dashboard - MailMemo</title>
</svelte:head>

<div class="mt-12 h-full w-full max-w-[750px] pb-16 lg:pb-0">
	<div class="flex w-full items-center justify-between">
		<div class="flex flex-col">
			<p class="font-header text-xl">Hello</p>
			<h3 class="font-header text-2xl">
				{$user?.name}
			</h3>
		</div>
		<a href="#">icon</a>
	</div>
	<div class="mt-3 flex gap-2">
		{#if data?.calendars?.total}
			{#each data.calendars.documents as calendar}
				<Badge variant="outline">
					<div class="h-2 w-2 rounded-full" style={`background-color: ${calendar.color}`}></div>
					<p class="ml-2">{calendar.name}</p>
				</Badge>
			{/each}
			<Badge variant="outline">
				+
				<p class="ml-2">Add calendar</p>
			</Badge>
		{/if}
	</div>
	<Separator class="mt-5" />
	<h1 class="font-header mt-6 text-xl">Today</h1>

	<div class="mt-4 grid gap-2">
		{#if todayEvents?.length}
			{#each todayEvents as event, i}
				{@const calendar = data.calendars.documents.find((c) => c.$id === event.calendarId)}
				<Card.Root class={`frosted ${evaluateCornerRadius(i, todayEvents?.length)}`}>
					<Card.Header class="p-4">
						<Card.Description>
							<div class="flex h-3 items-center gap-2">
								<p>
									{toLocaleTimeISO(event.startAt, false)} - {toLocaleTimeISO(event.endAt, false)}
								</p>
								<Separator orientation="vertical" class="h-full" />

								<span class="flex items-center gap-2">
									<span
										class="h-2 w-2 rounded-full"
										style={`background-color: ${calendar?.color}`}
									/>
									<p>{calendar?.name}</p>
								</span>
							</div>
						</Card.Description>

						<p class="mt-2">{event.name}</p>
					</Card.Header>
				</Card.Root>
			{/each}
		{:else}
			<EmptyCard />
		{/if}
	</div>
	<h1 class="font-header mt-8 text-xl">Upcoming Events</h1>

	<div class="mt-4 grid gap-2">
		{#if futureEvents?.length}
			{@const groupedByDay = Object.entries(futureEventsGrouped)}
			{#each groupedByDay as group, groupIndex}
				<Card.Root class={`frosted ${evaluateCornerRadius(groupIndex, groupedByDay?.length)}`}>
					<Card.Header class="p-4">
						<span class="mb-2">
							<Badge>
								In {group[0]} days
							</Badge>
						</span>
						{#each group[1] as event, i}
							{@const calendar = data.calendars.documents.find((c) => c.$id === event.calendarId)}
							<Separator class="mt-3" />
							<Card.Description class="pt-4">
								<div class="flex h-3 items-center gap-2">
									<p>
										{toLocaleTimeISO(event.startAt, false)} - {toLocaleTimeISO(event.endAt, false)}
									</p>
									<Separator orientation="vertical" class="h-full" />

									<span class="flex items-center gap-2">
										<span
											class="h-2 w-2 rounded-full"
											style={`background-color: ${calendar?.color}`}
										/>
										<p>{calendar?.name}</p>
									</span>
								</div>
							</Card.Description>

							<p class="mt-2" class:pb-4={i + 1 !== group[1]?.length}>{event.name}</p>
						{/each}
					</Card.Header>
				</Card.Root>
			{/each}
		{:else}
			<EmptyCard />
		{/if}
	</div>
</div>

<div
	class="fixed bottom-0 flex w-full max-w-[750px] justify-between gap-4 px-5 py-8 lg:relative lg:px-0"
>
	<Button variant="outline" class="w-full lg:w-auto">Sync calendars</Button>
	<Button class="w-full lg:w-auto">Send Email</Button>
</div>
