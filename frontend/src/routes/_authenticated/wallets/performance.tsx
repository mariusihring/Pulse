import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/wallets/performance')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_authenticated/wallets/performance"!</div>
}