<template>
  <div class="wallet-update">
    <h2>Wallet Update Test</h2>
    <input
      v-model="walletAddress"
      type="text"
      placeholder="Enter wallet address"
    />
    <button @click="startUpdate"> {{isFetching ? "Start Update" : "Loading..."}}</button>

    <div v-if="jobId">
      <p><strong>Job ID:</strong> {{ jobId }}</p>
      <p><strong>Subscription Data:</strong></p>
      <pre>{{ subscriptionData }}</pre>
    </div>
  </div>
</template>

<script  setup lang="ts">
import { ref, watch } from 'vue'
import { useMutation, useSubscription } from 'villus'
import {graphql} from "../src/gql"

// GraphQL mutation to start the wallet update process.
const START_WALLET_UPDATE = graphql(`
  mutation StartWalletUpdate($walletAddress: String!) {
    startWalletUpdate(walletAddress: $walletAddress) {
      id
      walletAddress
    }
  }
`)

// GraphQL subscription to listen for wallet updates.
const WALLET_UPDATES = graphql(`
  subscription WalletUpdates ($jobID: ID!) {
    walletUpdates(jobID: $jobID) {
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
`)
const walletAddress = ref('')
const jobId = ref('')
const subscriptionData = ref('')

// Create the mutation hook.
const { execute } = useMutation(START_WALLET_UPDATE)

// Create a reactive computed variable for subscription variables.
const subscriptionVariables = computed(() => ({
  jobID: jobId.value
}))

// Create a paused flag to control subscription execution.
const paused = ref(true)
const { subscribe, isFetching } = useSubscription({
  query: WALLET_UPDATES,
  paused,
  variables: subscriptionVariables
}, ({ data, error }) => {
  console.log(data)
  subscriptionData.value = JSON.stringify(data?.walletUpdates, null, 2)
})

const startUpdate = async () => {
  try {
    // Call the mutation to start the wallet update.
    const result = await execute({ walletAddress: walletAddress.value })
    if (result.data && !result.error) {
      console.log(result)
      jobId.value = result.data.startWalletUpdate.id as string
      paused.value = false
      subscribe()
    }
  } catch (error) {
    console.error('Error starting wallet update:', error)
  }
}
</script>

<style scoped>
.wallet-update {
  max-width: 600px;
  margin: 2rem auto;
  padding: 1rem;
  border: 1px solid #ccc;
  border-radius: 8px;
}
input {
  padding: 0.5rem;
  margin-right: 1rem;
  width: 60%;
}
button {
  padding: 0.5rem 1rem;
}
pre {
  background: #f4f4f4;
  padding: 1rem;
  border-radius: 4px;
}
</style>
