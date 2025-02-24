<script setup lang="ts">
import {useQuery, useSubscription} from "villus"

const GetWallet = `
  subscription GetWallet($address: String!) {
    walletUpdates (walletAddress: $address) {
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
    }
  }
  `

 
const obj = ref({})
useSubscription({
  query: GetWallet,
  variables: {address: "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3"},
}, ({data, error}) => {
  obj.value = data
  if (error) {
    console.error(error)
  }
})
</script>

<template>
  <div>
    <h1>Nuxt Routing set up successfully!</h1>
    <div>{{ JSON.stringify(obj.value) }}</div>

   <a href="https://nuxt.com/docs/getting-started/routing" target="_blank">Learn more about Nuxt Routing</a>
  </div>
</template>
