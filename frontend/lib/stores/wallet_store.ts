import {createStore} from "zustand/vanilla"
import {type WalletUpdate} from "@/lib/gql/graphql"
import {immer} from "zustand/middleware/immer"


export type WalletStoreState = {
    wallets: WalletUpdate[]
}


export type WalletStoreActions = {
    updateWallet: (update: WalletUpdate) => Promise<void>
}

export type WalletStore = WalletStoreState & WalletStoreActions

export const initialStoreState: WalletStoreState = {
    wallets: [],
}

export const createWalletStore = (initState: WalletStoreState = initialStoreState) => {
    return createStore<WalletStore>(
        immer((set) => ({
            ...initState,
            updateWallet: (update: WalletUpdate) => set((state) => {
                const index = state.wallets.findIndex(wallet => wallet.JobID === update.JobID)
                if (index !== -1) {
                    state.wallets[index] = update
                } else {
                    state.wallets.push(update)
                }
            }),
        }))
    )
}