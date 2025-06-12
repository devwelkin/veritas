import type {
	ErrorResponse,
	ShortenUrlRequest,
	ShortenUrlResponse,
} from "@/types/api";

const API_BASE_URL =
	import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export const shortenUrl = async (
	request: ShortenUrlRequest,
): Promise<ShortenUrlResponse> => {
	const response = await fetch(`${API_BASE_URL}/create`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		const errorData: ErrorResponse = await response.json();
		throw new Error(
			errorData.message || "Something went wrong during shortening.",
		);
	}

	return response.json();
}; 