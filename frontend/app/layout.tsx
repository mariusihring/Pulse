'use client'
import "./globals.css";
import {Provider, cacheExchange, Client, fetchExchange, subscriptionExchange } from "urql";
import {createClient as createWSClient} from "graphql-ws"
import { WalletStoreProvider} from "@/lib/providers/wallet_provider";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {

  const wsClient = createWSClient({
    url: "ws://localhost:3001/query"
  })

  const client = new Client({
    url: "http://localhost:3001/query",
    exchanges: [cacheExchange, fetchExchange, subscriptionExchange({forwardSubscription(request) {
      const input = {...request, query: request.query || ""};
      return {
        subscribe(sink) {
          const unsubscribe = wsClient.subscribe(input, sink);
          return {unsubscribe}
        }
      }
    }})],
    
  })
  return (
    <html lang="en">
      <body
      >
        <Provider value={client}>
          <WalletStoreProvider>

            {children}
          </WalletStoreProvider>
        </Provider>
      </body>
    </html>
  );
}
