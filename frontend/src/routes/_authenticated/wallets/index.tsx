import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/wallets/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_authenticated/wallets/"!</div>
}
