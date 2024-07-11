import { writable, type Writable } from 'svelte/store';

export const step = writable(1);
export const preferences: Writable<{
	frequency: string | null;
	hour: string;
	minute: string;
	format: string;
	frequencyDetails: string | null;
	period: string | null;
	timezone: string | null;
}> = writable({
	frequency: null,
	hour: '09',
	minute: '00',
	format: 'AM',
	frequencyDetails: null,
	period: null,
	timezone: null
});
