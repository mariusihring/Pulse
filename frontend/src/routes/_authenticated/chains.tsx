import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/chains')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_authenticated/chains"!</div>
}
