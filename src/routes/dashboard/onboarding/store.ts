import { writable, type Writable } from 'svelte/store';

export const step = writable(1);
export const preferences: Writable<{
	frequency: string;
	hour: string;
	minute: string;
	format: string;
	frequencyDetails: string;
	period: string;
	timezone: string;
}> = writable({});
