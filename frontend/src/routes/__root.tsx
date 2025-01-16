import * as React from "react";
import {
	Outlet,
	createRootRouteWithContext,
} from "@tanstack/react-router";
import { AuthProvider, useAuth } from "@/contexts/auth_context";
import type { QueryClient } from "@tanstack/react-query";

interface RouterContext {
	queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
	component: () => (
		<>
			<AuthProvider>
				<React.Fragment>
					<Outlet />
				</React.Fragment>
			</AuthProvider>
		</>
	),
});
