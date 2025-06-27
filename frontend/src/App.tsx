import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { AnimatePresence, motion } from "framer-motion";
import { Check, Copy, Loader2, AlertTriangle } from "lucide-react";
import { shortenUrl as shortenUrlApi } from "@/api/urlApi";

function App() {
	const [longUrl, setLongUrl] = useState("");
	const [shortUrl, setShortUrl] = useState("");
	const [loading, setLoading] = useState(false);
	const [copied, setCopied] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		if (!longUrl) return;

		setLoading(true);
		setError(null);
		try {
			const data = await shortenUrlApi({ original_url: longUrl });
			if (data.short_url) {
				const fullShortUrl = `${window.location.protocol}//${window.location.host}/${data.short_url}`;
				setShortUrl(fullShortUrl);
			}
		} catch (err: unknown) {
			if (err instanceof Error) {
				setError(err.message || "An unexpected error occurred.");
			} else {
				setError("An unexpected error occurred.");
			}
		} finally {
			setLoading(false);
		}
	};

	const handleCopy = () => {
		navigator.clipboard.writeText(shortUrl);
		setCopied(true);
		setTimeout(() => setCopied(false), 2000);
	};

	const handleReset = () => {
		setLongUrl("");
		setShortUrl("");
		setError(null);
	};

	return (
		<main className="dark font-sans flex min-h-screen w-full flex-col items-center justify-center bg-slate-900 text-white">
			<div className="w-full max-w-md p-4">
				<AnimatePresence mode="wait">
					{!shortUrl ? (
						<motion.div
							key="form"
							initial={{ opacity: 0, y: 20 }}
							animate={{ opacity: 1, y: 0 }}
							exit={{ opacity: 0, y: -20 }}
							transition={{ duration: 0.3 }}
							className="space-y-6"
						>
							<div className="text-center">
								<h1 className="text-4xl font-bold tracking-tight">
									Tame your links
								</h1>
								<p className="text-slate-400">
									Short, sweet, and to the point.
								</p>
							</div>
							<form onSubmit={handleSubmit} className="flex flex-col gap-4">
								<Input
									type="url"
									placeholder="Paste your long URL here..."
									value={longUrl}
									onChange={(e) => setLongUrl(e.target.value)}
									className="h-12 bg-slate-800 border-slate-700 text-lg placeholder:text-slate-500 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:ring-offset-slate-900"
									required
								/>
								<Button
									type="submit"
									className="h-12 text-lg font-bold text-white bg-gradient-to-r from-violet-500 to-indigo-600 hover:opacity-90 transition-opacity disabled:opacity-50"
									disabled={loading || !longUrl}
								>
									{loading ? (
										<Loader2 className="animate-spin" />
									) : (
										"Shorten"
									)}
								</Button>
							</form>
							{error && (
								<div className="flex items-center gap-2 rounded-md border border-red-500/50 bg-red-500/10 p-3 text-red-400">
									<AlertTriangle className="h-5 w-5" />
									<p className="text-sm">{error}</p>
								</div>
							)}
						</motion.div>
					) : (
						<motion.div
							key="result"
							initial={{ opacity: 0, y: 20 }}
							animate={{ opacity: 1, y: 0 }}
							exit={{ opacity: 0, y: -20 }}
							transition={{ duration: 0.3 }}
						>
							<Card className="bg-slate-800 border-slate-700">
								<CardHeader>
									<CardTitle className="text-2xl">Your short link is ready!</CardTitle>
									<CardDescription className="text-slate-400">
										Share it, save it, use it.
									</CardDescription>
								</CardHeader>
								<CardContent className="space-y-4">
									<div className="flex items-center gap-2 rounded-md border border-slate-700 bg-slate-900 p-3">
										<p className="font-mono text-lg text-indigo-400 overflow-x-auto whitespace-nowrap">
											{shortUrl}
										</p>
										<Button
											variant="ghost"
											size="icon"
											onClick={handleCopy}
											className="ml-auto flex-shrink-0 text-slate-400 hover:bg-slate-700 hover:text-white"
										>
											{copied ? (
												<Check className="h-5 w-5 text-green-500" />
											) : (
												<Copy className="h-5 w-5" />
											)}
										</Button>
									</div>
									<Button
										onClick={handleReset}
										className="w-full h-12 text-lg font-bold text-white bg-gradient-to-r from-violet-500 to-indigo-600 hover:opacity-90 transition-opacity"
									>
										Shorten another one
									</Button>
								</CardContent>
							</Card>
						</motion.div>
					)}
				</AnimatePresence>
			</div>
		</main>
	);
}

export default App; 