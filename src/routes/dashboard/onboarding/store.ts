import { writable } from 'svelte/store';

export const step = writable(1);
export const preferences = writable({});
