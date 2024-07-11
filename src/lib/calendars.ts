export const CalColors = ['pink', 'orange', 'purple', 'mint', 'blue', 'yellow'] as const;

export type Calendar = {
	name: string;
	url: string;
	color: typeof CalColors;
};

export type CalEvent = {
	$id: string;
	$permissions: string[];
	$createdAt: string;
	$collectionId: string;
	$updatedAt: string;
	$databaseId: string;
	name: string;
	startAt: string;
	endAt: string;
	calendarId: string;
	modifiedAt: string;
	uid: string;
};

export function colorToHex(color: (typeof CalColors)[number]) {
	switch (color) {
		case 'pink':
			return '#FD366E';
		case 'orange':
			return '#FE9567';
		case 'purple':
			return '#7C67FE';
		case 'mint':
			return '#85DBD8';
		case 'blue':
			return '#68A3FE';
		case 'yellow':
			return '#FED367';
	}
}
