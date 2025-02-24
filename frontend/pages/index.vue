<template>
  <div class="wallet-update">
    <h2>Wallet Update Test</h2>
    <input
      v-model="walletAddress"
      type="text"
      placeholder="Enter wallet address"
    />
    <Button @click="startUpdate"> {{isFetching ? "Start Update" : "Loading..."}}</Button>

    <div v-if="jobId">
      <p><strong>Job ID:</strong> {{ jobId }}</p>
    </div>
  <div class="flex w-full items-center justify-center">
    <UserTable :columns="columns" :data="tokens"/>
  </div>
  </div>
  

</template>

<script  setup lang="ts">
import { ref, watch } from 'vue'
import { useMutation, useSubscription } from 'villus'
import {graphql} from "../src/gql"
import UserTable from "../src/compontents/test_table.vue"
import type { WalletUpdate } from '@/src/gql/graphql'
import {Button} from "@/components/ui/button"
import {Checkbox} from "@/components/ui/checkbox"
import { computed } from 'vue'
import {createColumnHelper} from "@tanstack/vue-table";
import type {Token} from "@/src/gql/graphql"

const columnHelper = createColumnHelper<Token>()
const columns = [
    columnHelper.display({
      id: "select",
      header: ({table}) => h(Checkbox, {
        "modelValue": table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && "indeterminate"),
        'onUpdate:modelValue': value => table.toggleAllPageRowsSelected(!!value),
        'ariaLabel':  "Select all",
      }),
      cell: ({row}) => {
        return h(Checkbox, {
          'modelValue': row.getIsSelected(),
          'onUpdate:modelValue': value => row.toggleSelected(!!value),
          'ariaLabel':  "Select row",
        })
      }
    }),
    columnHelper.display({
      id: "name",
      header: ({table}) => h("h1", {

      }, "Name"),
      cell: ({row}) => {
        console.log(row.getValue('name'))
        return h('p', {}, row.getValue('name'))
      }
    })
]

const tokens = computed(() => {
  // Re-create the array if it exists, to trigger reactivity.
  return subscriptionData.value?.Wallet?.tokens ? [ ...subscriptionData.value.Wallet.tokens ] : []
})

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
const subscriptionData = ref<WalletUpdate | undefined>(undefined)

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
  subscriptionData.value = data?.walletUpdates
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
