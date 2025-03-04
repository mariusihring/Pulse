'use client'

import { createContext, useContext, useRef, ReactNode } from 'react'
import { useStore } from 'zustand'
import type { WalletStore } from '@/lib/stores/wallet_store'
import { createWalletStore } from '@/lib/stores/wallet_store'

// Provide a typed context
const CounterStoreContext = createContext<ReturnType<typeof createWalletStore> | null>(null)

interface CounterStoreProviderProps {
    children: ReactNode
}

export const CounterStoreProvider = ({ children }: CounterStoreProviderProps) => {
    const storeRef = useRef<ReturnType<typeof createWalletStore>>()
    if (!storeRef.current) {
        storeRef.current = createWalletStore()
    }

    return (
        <CounterStoreContext.Provider value={storeRef.current}>
            {children}
        </CounterStoreContext.Provider>
    )
}

// Generic hook to allow proper typing of the selector function
export const useCounterStore = <T,>(
    selector: (state: WalletStore) => T,
    equalityFn?: (a: T, b: T) => boolean
): T => {
    const store = useContext(CounterStoreContext)
    if (!store) {
        throw new Error('useCounterStore must be used within a CounterStoreProvider')
    }
    return useStore<WalletStore, T>(store, selector, equalityFn)
}
