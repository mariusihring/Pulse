<template>
  <div>
    <NuxtRouteAnnouncer />
    
  </div>
  <NuxtPage />
</template>


<script setup lang="ts">
import {defaultPlugins, handleSubscriptions, useClient} from "villus"
import {createClient} from "graphql-ws"
const wsClient = createClient({
  url: "ws://localhost:3001/query"
})
const subscriptionHandler = (handleSubscriptions(operation => {
  return {
    subscribe: obs => {
      wsClient.subscribe(
        {
          query: operation.query,
          variables: operation.variables
        },
        obs,
      );

      return {
        unsubscribe: () => {}
      }
    }
  }
}))
useClient({
  url: "http://localhost:3001/query",
  use: [...defaultPlugins()]
})

</script>
