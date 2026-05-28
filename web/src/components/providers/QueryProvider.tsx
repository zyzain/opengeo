"use client";

import "@ant-design/v5-patch-for-react-19";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { App } from "antd";
import { useState } from "react";

export default function QueryProvider({
	children,
}: { children: React.ReactNode }) {
	const [queryClient] = useState(
		() =>
			new QueryClient({
				defaultOptions: {
					queries: {
						staleTime: 60 * 1000,
						retry: 1,
						refetchOnWindowFocus: false,
					},
				},
			}),
	);

	return (
		<App>
			<QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
		</App>
	);
}
