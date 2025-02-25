'use client'
import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import {Provider, cacheExchange, Client, fetchExchange, subscriptionExchange } from "urql";
import {createClient as createWSClient} from "graphql-ws"

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
        {children}
        </Provider>
      </body>
    </html>
  );
}
