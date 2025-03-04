"use client"

import type { ReactNode } from "react"
import { useTheme } from "next-themes"
import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useWalletStore} from "@/lib/providers/wallet_provider";
import {graphql} from "@/lib/gql";
import {useMutation, useSubscription} from "urql";

interface LayoutProps {
  children: ReactNode
}


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
        name
        description
        network
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
          history_prices {
            Timestamp
            Open
            High
            Low
            Close
            Volume
          }
        }
        transactions {
          jsonrpc
          id
          result {
            block_time
            slot
            meta {
              compute_units_consumed
              fee
              log_messages
              post_balances
              pre_balances
              post_token_balances {
                account_index
                mint
                owner
                program_id
                ui_token_amount {
                  amount
                  decimals
                  ui_amount
                  ui_amount_string
                }
              }
              inner_instructions {
                index
                instructions {
                  accounts
                  data
                  program_id_index
                  stack_height
                }
              }
              pre_token_balances {
                account_index
                mint
                owner
                program_id
                ui_token_amount {
                  amount
                  decimals
                  ui_amount
                  ui_amount_string
                }
              }
              rewards {
                info
              }
              status {
                ok
                error_message
              }
            }
            transaction {
              signatures
              message {
                account_keys
                recent_blockhash
                address_table_lookups {
                  account_key
                  readonly_indexes
                  writable_indexes
                }
                header {
                  num_readonly_signed_accounts
                  num_readonly_unsigned_accounts
                  num_required_signatures
                }
                instructions {
                  accounts
                  data
                  program_id_index
                  stack_height
                }
              }
            }
          }
          err {
            code
            message
          }
        }
      }
    }
  }
`);
export default function Layout({ children }: LayoutProps) {
  const { theme } = useTheme()
  const [mounted, setMounted] = useState(false)
  const {updateWallet} = useWalletStore((state) => state)
  const [address, setAddress] = useState("");
  const [loading, setLoading] = useState(false)
  const [jobId, setJobId] = useState("");
  const [curr_progress, setCurrProcess] = useState(0);
  const [responseData, setResponseData] = useState("");
  const [startWalletScanResult, startWalletScan] = useMutation(StartWalletScanMutation);

  // Define the subscription with pause initially true
  //@ts-expect-error
  const [subscriptionResult] = useSubscription({
    query: WalletUpdateSubscription,
    variables: { input: jobId },
    pause: !jobId, // Only start subscription when we have a jobId
  });

  // Update the textarea whenever new subscription data comes in
  useEffect(() => {
    if (subscriptionResult.data) {
      if (subscriptionResult.data.walletUpdates.Progress === 100) {
        setLoading(false)
      }
      // Format the JSON with 2-space indentation for readability
      const formattedData = JSON.stringify(subscriptionResult.data, null, 2);
      updateWallet(subscriptionResult.data.walletUpdates)
      setCurrProcess(subscriptionResult.data.walletUpdates.Progress)
      setResponseData(formattedData);
    }
  }, [subscriptionResult.data]);

  const handleFetch = async () => {
    // Reset response data when starting a new fetch
    setResponseData("");
    setLoading(true)

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

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return null
  }

  return (
    <div className={`flex h-screen ${theme === "dark" ? "dark" : ""}`}>
      <div className="w-full flex flex-1 flex-col">
        <header className="h-16 border-b border-gray-200 dark:border-[#1F1F23] flex gap-2 w-full items-center justify-end p-2 px-6 ">
          <Input
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Insert wallet address"
              className="flex-grow"
              required
          />
          <Button onClick={handleFetch} disabled={startWalletScanResult.fetching}>
            {loading ? "Loading..." : "Fetch"}
          </Button>
        </header>
        <main className="flex-1 overflow-auto p-6 bg-white dark:bg-[#0F0F12]">{children}</main>
      </div>
    </div>
  )
}

