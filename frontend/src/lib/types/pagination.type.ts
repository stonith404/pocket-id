export type PaginationRequest = {
	page: number;
	limit: number;
};

export type PaginationResponse = {
	totalPages: number;
	totalItems: number;
	currentPage: number;
};

export type Paginated<T> = {
	data: T[];
	pagination: PaginationResponse;
};