import {createStore} from "zustand/vanilla"
import {type WalletUpdate} from "@/lib/gql/graphql"
import {immer} from "zustand/middleware/immer"


export type WalletStoreState = {
    wallets: WalletUpdate[]
}


export type WalletStoreActions = {
}

export type WalletStore = WalletStoreState & WalletStoreActions

export const initState: WalletStoreState = {
    wallets: [],
}

export const createWalletStore = (initState: WalletStoreState = initState) => {
    return createStore<WalletStore>(
        immer((set) => ({
            ...initState,
            //decrement: () => set((state) => { state.count-- }),
            //increment: () => set((state) => { state.count++ }),
        }))
    )
}