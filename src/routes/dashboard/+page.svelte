<script lang="ts">
	import { user } from '$lib/stores';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
	import * as Card from '$lib/components/ui/card';
	import { toLocaleTimeISO } from '$lib/utils';
	import EmptyCard from './emptyCard.svelte';

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
	$: console.log(data);
	$: todayEvents = data?.events?.documents?.filter((event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		return eventDate.getDate() === today.getDate();
	});
</script>

<div class=" mt-12 h-full w-full">
	<div class="flex w-full items-center justify-between">
		<div class="flex flex-col">
			<p class="text-xl">Hello</p>
			<h3 class="text-2xl font-bold">
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
	<h1 class="mt-6 text-xl">Today</h1>

	<div class="mt-4 grid gap-1">
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
	<h1 class="mt-8 text-xl">Upcoming Events</h1>

	<div class="mt-4 grid gap-4">
		{#if data?.events?.total}
			{#each data.events.documents as event}
				<div class="grid gap-2">
					<p class="text-sm text-muted-foreground">{event.start}</p>
					<p class="text-lg">{event.title}</p>
				</div>
			{/each}
		{:else}
			<EmptyCard />
		{/if}
	</div>
</div>
