import { writable, type Writable } from 'svelte/store';

export const step = writable(1);
export const preferences: Writable<{
	frequency: string | null;
	hour: string | null;
	minute: string | null;
	format: string | null;
	frequencyDetails: string | null;
	period: string | null;
	timezone: string | null;
}> = writable({
	frequency: null,
	hour: null,
	minute: null,
	format: null,
	frequencyDetails: null,
	period: null,
	timezone: null
});
