'use client'
import { useState, useEffect } from "react";
import { useMutation, useSubscription } from "urql";
import { graphql } from "@/lib/gql";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {Progress} from "@/components/ui/progress";

const StartWalletScanMutation = graphql(`
    mutation StartWalletScan($input: String!) {
        startWalletUpdate(walletAddress: $input) {
            id
            walletAddress
        }
    }
`);

const WalletUpdateSubscription = graphql(`
    subscription WalletUpdate($input: ID!) {
        walletUpdates(jobID: $input) {
            JobID
            Progress
            Wallet {
                address
                sol_balance
                sol_value
                wallet_value
                last_updated
                tokens {
                    name
                    address
                    pool
                    description
                    image
                    amount
                    price
                    pnl
                    invested
                    value
                    history_prices
                }
                transactions {
                    jsonrpc
                    id
                }
            }
        }
    }
`);

export default function Home() {
    const [address, setAddress] = useState("");
    const [jobId, setJobId] = useState("");
    const [curr_progress, setCurrProcess] = useState(0);
    const [responseData, setResponseData] = useState("");
    const [startWalletScanResult, startWalletScan] = useMutation(StartWalletScanMutation);

    // Define the subscription with pause initially true
    const [subscriptionResult] = useSubscription({
        query: WalletUpdateSubscription,
        variables: { input: jobId },
        pause: !jobId, // Only start subscription when we have a jobId
    });

    // Update the textarea whenever new subscription data comes in
    useEffect(() => {
        if (subscriptionResult.data) {
            // Format the JSON with 2-space indentation for readability
            const formattedData = JSON.stringify(subscriptionResult.data, null, 2);
            setCurrProcess(subscriptionResult.data.walletUpdates.Progress)
            setResponseData(formattedData);
        }
    }, [subscriptionResult.data]);

    const handleFetch = async () => {
        // Reset response data when starting a new fetch
        setResponseData("");

        // Call the mutation with the wallet address
        const result = await startWalletScan({ input: address });
        if (result.error) {
            console.error("Mutation Error:", result.error);
            setResponseData(JSON.stringify({ error: result.error }, null, 2));
            return;
        }

        // Extract the id from the mutation response
        const newJobId = result.data?.startWalletUpdate?.id;
        if (newJobId) {
            // Update the jobId state which will automatically start the subscription
            setJobId(newJobId);

            // Show initial response
            setResponseData(JSON.stringify({
                message: "Subscription started",
                jobId: newJobId
            }, null, 2));
        } else {
            console.error("No job id returned from mutation");
            setResponseData(JSON.stringify({
                error: "No job ID returned from mutation"
            }, null, 2));
        }
    };

    return (
        <div className="p-5 w-full flex flex-col gap-4">
            <div className="flex gap-4">
                <Input
                    value={address}
                    onChange={(e) => setAddress(e.target.value)}
                    placeholder="Insert wallet address"
                    className="flex-grow"
                />
                <Button onClick={handleFetch} disabled={startWalletScanResult.fetching}>
                    {subscriptionResult.fetching ? "Loading..." : "Fetch"}
                </Button>
            </div>

            {jobId && (
                <div className="text-sm text-gray-500">
                    Job ID: {jobId} {subscriptionResult.fetching && "• Receiving updates..."}
                </div>
            )}
            <Progress value={curr_progress}/>
            <Textarea
                value={responseData}
                readOnly
                className="min-h-64 font-mono text-sm"
                placeholder="Subscription data will appear here..."
            />

            {subscriptionResult.error && (
                <div className="text-red-500">
                    Error: {subscriptionResult.error.message}
                </div>
            )}
        </div>
    );
}