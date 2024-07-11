export type Calendar = {
	name: string;
	url: string;
	color: string;
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
