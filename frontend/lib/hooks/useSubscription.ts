import {useQuery, useQueryClient} from "@tanstack/react-query"
import {request, gql} from "graphql-request"
import { useEffect } from "react"
import {createClient} from "graphql-ws"

const endpoint = "http://localhost:3001/query"

const wsClient = createClient({
    url: 'wss://127.0.0.1:3001/query'
})


export function useSubscription(queryKey, queryFn, subscriptionQuery) {
    const queryClient = useQueryClient()

    const queryResult = useQuery({
        queryKey,
        queryFn
    })

    useEffect(() => {
        if (!subscriptionQuery) return;

        const unsubscribe = wsClient.subscribe(
            {query:  subscriptionQuery},
            {next: (data) => {
                queryClient.setQueryData(queryKey, (oldData) => {
                    return data
                })
            },
            error: (error) => {
                console.error("Subscription error: ", error)
            },
            complete: () => {
                console.log("SubscriptionComplete")
            }
        
        }
        );

        return () => {
            unsubscribe()
        }
    }, [queryClient, queryKey, subscriptionQuery])
    return queryResult
}


