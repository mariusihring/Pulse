import * as React from "react";
import { Outlet, createRootRoute, useNavigate } from "@tanstack/react-router";
import { AuthProvider, useAuth } from "@/contexts/auth_context";

export const Route = createRootRoute({
  component: RootComponent,
 
});

function RootComponent() {
  return (
    <AuthProvider>
      <React.Fragment>
      <Outlet />
    </React.Fragment>
</AuthProvider>
    
  );
}
