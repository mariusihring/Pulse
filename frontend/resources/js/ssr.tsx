import { createInertiaApp } from '@inertiajs/react';
import createServer from '@inertiajs/react/server';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import ReactDOMServer from 'react-dom/server';
import { type RouteName, route } from 'ziggy-js';
import {Toaster} from  "@/components/ui/sonner"
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: 1, // Retry failed queries once
        staleTime: 1000 * 60 * 5, // Cache data for 5 minutes
      },
    },
  });
const appName = import.meta.env.VITE_APP_NAME || 'Laravel';

createServer((page) =>
    createInertiaApp({
        page,
        render: ReactDOMServer.renderToString,
        title: (title) => `${title} - ${appName}`,
        resolve: (name) => resolvePageComponent(`./pages/${name}.tsx`, import.meta.glob('./pages/**/*.tsx')),
        setup: ({ App, props }) => {
            /* eslint-disable */
            // @ts-expect-error
            global.route<RouteName> = (name, params, absolute) =>
                route(name, params as any, absolute, {
                    // @ts-expect-error
                    ...page.props.ziggy,
                    // @ts-expect-error
                    location: new URL(page.props.ziggy.location),
                });
            /* eslint-enable */

            return
            <QueryClientProvider client={queryClient}>
            <App {...props} />
            <Toaster richColors/>
            <ReactQueryDevtools initialIsOpen={false} />
            </QueryClientProvider>
        },
    }),
);
