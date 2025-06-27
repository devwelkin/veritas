import type {
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

	const data = await response.json();

	if (!response.ok || data.error) {
		throw new Error(data.error || "Something went wrong during shortening.");
	}

	return data;
}; 