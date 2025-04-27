export interface Token {
    id: string;
    chain_id: string;
    name: string;
    current_price: string;
    logo: string | null;
    symbol: string;
    address: string;
    mint: string;
    created_at: string;
    updated_at: string;
}

export interface TokenHolding {
    id: string;
    user_id: string;
    wallet_id: string;
    token_id: string;
    amount: string;
    value: string;
    created_at: string;
    updated_at: string;
    token: Token;
}

export interface Snapshot {
    id: string;
    wallet_id: string;
    value: string;
    created_at: string;
    updated_at: string;
}

export interface TokenSwap {
    id: string;
    chain_id: string;
    token_id: string;
    wallet_id: string;
    transaction_hash: string;
    transaction_type: string;
    transaction_index: number;
    sub_category: string;
    block_timestamp: string;
    block_number: number;
    wallet_address: string;
    pair_address: string;
    pair_label: string;
    exchange_address: string;
    exchange_name: string;
    exchange_logo: string;
    base_token: string;
    quote_token: string;
    bought: {
        address: string;
        amount: string;
        usdPrice: number;
        usdAmount: number;
        symbol: string;
        logo: string | null;
        name: string;
        tokenType: string;
    };
    sold: {
        address: string;
        amount: string;
        usdPrice: number;
        usdAmount: number;
        symbol: string;
        logo: string | null;
        name: string;
        tokenType: string;
    };
    base_quote_price: string;
    total_value_usd: string;
}

export interface Wallet {
    id: string;
    address: string;
    name: string;
    chain_token_amount: number;
    value: number;
    chain_id: string;
    user_id: string;
    favorite: boolean;
    created_at: string;
    updated_at: string;
    snapshots: Snapshot[];
    tokenswaps: TokenSwap[];
}

export interface User {
    id: string;
    name: string;
    email: string;
    email_verified_at: string | null;
    created_at: string;
    updated_at: string;
    wallets: Wallet[];
    token_holdings: TokenHolding[];
}
