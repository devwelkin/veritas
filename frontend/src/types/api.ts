export interface ShortenUrlRequest {
  original_url: string;
}

export interface ShortenUrlResponse {
  short_url: string;
}

export interface ErrorResponse {
  error: string;
} 