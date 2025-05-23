


export enum TransactionType {
    Buy = "buy",
    Sell = "sell",
}

export enum SubCategory {
    NewPosition = "newPosition",
    Accumulation = "accumulation",
    SellAll = "sellAll",
    PartialSell = "partialSell",
}

export enum TokenType {
    Token0 = "token0",
    Token1 = "token1",
}

export interface Token {
    id: string; // UUID
    chain_id: string; // UUID
    name: string;
    current_price: string; // Decimal as string
    logo: string | null;
    symbol: string;
    address: string;
    mint: string;
    pnl: TokenPnl;
    created_at: string; // ISO 8601 timestamp
    updated_at: string; // ISO 8601 timestamp
}

export interface TokenPnl {
    averageBuyPrice: number;
    realizedPnl: number;
    remainingAmount: number;
    totalBoughtUsd: number;
    totalPnl: number;
    totalSoldUsd: number;
    totalTokenAmountBought: number;
    totalTokenAmountSold: number;
    unrealizedPnl: number;
}

export interface TokenHolding {
    id: string; // UUID
    user_id: string; // UUID
    wallet_id: string; // UUID
    token_id: string; // UUID
    amount: string; // Decimal as string
    value: string; // Decimal as string
    created_at: string;
    updated_at: string;
    token: Token;
}

export interface BoughtOrSoldToken {
    address: string;
    amount: string; // Decimal as string
    usdPrice: number;
    usdAmount: number;
    symbol: string;
    logo: string | null;
    name: string;
    tokenType: TokenType;
}

export interface TokenSwap {
    id: string; // UUID
    chain_id: string; // UUID
    token_id: string; // UUID
    wallet_id: string; // UUID
    transaction_hash: string;
    transaction_type: TransactionType;
    transaction_index: number;
    sub_category: SubCategory;
    block_timestamp: string; // ISO 8601 timestamp
    block_number: number;
    wallet_address: string;
    pair_address: string;
    pair_label: string;
    exchange_address: string;
    exchange_name: string;
    exchange_logo: string;
    base_token: string;
    quote_token: string;
    bought: BoughtOrSoldToken;
    sold: BoughtOrSoldToken;
    base_quote_price: string; // Decimal as string
    total_value_usd: string; // Decimal as string
}

export interface Snapshot {
    id: string; // UUID
    wallet_id: string; // UUID
    value: string; // Decimal as string
    created_at: string;
    updated_at: string;
}

export interface Chain {
    id: string; // UUID
    name: string;
    created_at: string;
    updated_at: string;
}

export interface Wallet {
    id: string; // UUID
    address: string;
    name: string;
    chain_token_amount: number;
    value: number;
    chain_id: string; // UUID
    user_id: string; // UUID
    favorite: boolean;
    created_at: string;
    updated_at: string;
    chain: Chain;
    snapshots: Snapshot[];
    tokenswaps: TokenSwap[];
    token_holdings: TokenHolding[];
}

export interface User {
    id: string; // UUID
    name: string;
    email: string;
    tokens: Token[];
    email_verified_at: string | null;
    created_at: string;
    updated_at: string;
    wallets: Wallet[];
    token_holdings: TokenHolding[];
}


