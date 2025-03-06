import { createStore } from "zustand/vanilla";
import { type WalletUpdate } from "@/lib/gql/graphql";
import { immer } from "zustand/middleware/immer";

export type WalletStoreState = {
    wallets: WalletUpdate[];
};

export type WalletStoreActions = {
    updateWallet: (update: WalletUpdate) => Promise<void>;
};

export type WalletStore = WalletStoreState & WalletStoreActions;

export const initialStoreState: WalletStoreState = {
    wallets: [],
};

export const createWalletStore = (initState: WalletStoreState = initialStoreState) => {
    return createStore<WalletStore>(
        immer((set) => ({
            ...initState,
            updateWallet: async (update: WalletUpdate) =>
                set((state) => {
                    const deduplicateTransactions = (transactions: any[] = []) => {
                        return Array.from(
                            new Map(
                                transactions.map((tx) => [tx.result?.block_time, tx])
                            ).values()
                        );
                    };

                    const index = state.wallets.findIndex(
                        (wallet) => wallet.JobID === update.JobID
                    );

                    if (index !== -1) {
                        const existingTransactions = state.wallets[index].Wallet.transactions || [];
                        const newTransactions = update.Wallet.transactions || [];
                        const mergedTransactions = [...existingTransactions, ...newTransactions];

                        update.Wallet.transactions = deduplicateTransactions(mergedTransactions);
                        state.wallets[index] = update;
                    } else {
                        if (update.Wallet && update.Wallet.transactions) {
                            update.Wallet.transactions = deduplicateTransactions(update.Wallet.transactions);
                        }
                        state.wallets.push(update);
                    }
                }),
        }))
    );
};
