/* eslint-disable */
import { DocumentTypeDecoration } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Decimal: { input: any; output: any; }
  Time: { input: any; output: any; }
  UUID: { input: any; output: any; }
};

export type Alert = {
  __typename?: 'Alert';
  condition: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  notificationSettings?: Maybe<Scalars['String']['output']>;
  token: Token;
  updatedAt: Scalars['Time']['output'];
  user: User;
};

export type AuthResponse = {
  __typename?: 'AuthResponse';
  token: Scalars['String']['output'];
  user: User;
};

export type Chain = {
  __typename?: 'Chain';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  subwallets: Array<Subwallet>;
  updatedAt: Scalars['Time']['output'];
};

export type CreateAlertInput = {
  condition: Scalars['String']['input'];
  notificationSettings?: InputMaybe<Scalars['String']['input']>;
  tokenId: Scalars['UUID']['input'];
};

export type CreateUserInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type HistoricalPrice = {
  __typename?: 'HistoricalPrice';
  createdAt: Scalars['Time']['output'];
  date: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  price: Scalars['Decimal']['output'];
  token: Token;
};

export type LoginInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  assignRole: User;
  createAlert: Alert;
  createSubwallet: Subwallet;
  createWallet: Wallet;
  deleteAlert: Scalars['Boolean']['output'];
  login: AuthResponse;
  register: AuthResponse;
  removeRole: User;
};


export type MutationAssignRoleArgs = {
  roleName: Scalars['String']['input'];
  userId: Scalars['UUID']['input'];
};


export type MutationCreateAlertArgs = {
  input: CreateAlertInput;
};


export type MutationCreateSubwalletArgs = {
  input: CreateSubwalletInput;
};


export type MutationCreateWalletArgs = {
  input: CreateWalletInput;
};


export type MutationDeleteAlertArgs = {
  id: Scalars['UUID']['input'];
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationRemoveRoleArgs = {
  roleName: Scalars['String']['input'];
  userId: Scalars['UUID']['input'];
};

export type PortfolioMetric = {
  __typename?: 'PortfolioMetric';
  calculatedAt: Scalars['Time']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  metricName: Scalars['String']['output'];
  metricValue: Scalars['Decimal']['output'];
  user: User;
};

export type PortfolioStats = {
  __typename?: 'PortfolioStats';
  dailyChange: Scalars['Decimal']['output'];
  topPerformers: Array<SubwalletToken>;
  totalPnl: Scalars['Decimal']['output'];
  totalValue: Scalars['Decimal']['output'];
  worstPerformers: Array<SubwalletToken>;
};

export type Query = {
  __typename?: 'Query';
  alerts: Array<Alert>;
  chain: Chain;
  chains: Array<Chain>;
  me: User;
  portfolioMetrics: Array<PortfolioMetric>;
  portfolioStats: PortfolioStats;
  roles: Array<Role>;
  subwallet: Subwallet;
  token: Token;
  tokenPriceHistory: Array<HistoricalPrice>;
  tokens: Array<Token>;
  transactionCategories: Array<TransactionCategory>;
  transactions: Array<Transaction>;
  wallet: Wallet;
  wallets: Array<Maybe<Wallet>>;
};


export type QueryChainArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryPortfolioMetricsArgs = {
  from: Scalars['Time']['input'];
  to: Scalars['Time']['input'];
};


export type QuerySubwalletArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryTokenArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryTokenPriceHistoryArgs = {
  from: Scalars['Time']['input'];
  id: Scalars['UUID']['input'];
  to: Scalars['Time']['input'];
};


export type QueryWalletArgs = {
  id: Scalars['UUID']['input'];
};

export type RegisterInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Role = {
  __typename?: 'Role';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  users: Array<User>;
};

export type Snapshot = {
  __typename?: 'Snapshot';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  snapshotDate: Scalars['Time']['output'];
  subwallet: Subwallet;
  totalPnl: Scalars['Decimal']['output'];
  totalValue: Scalars['Decimal']['output'];
};

export type Subwallet = {
  __typename?: 'Subwallet';
  address: Scalars['String']['output'];
  chain: Chain;
  createdAt: Scalars['Time']['output'];
  currentValue: Scalars['Float']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  snapshots: Array<Maybe<Snapshot>>;
  tokens: Array<Maybe<SubwalletToken>>;
  updatedAt: Scalars['Time']['output'];
};

export type SubwalletToken = {
  __typename?: 'SubwalletToken';
  amount: Scalars['Decimal']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  snapshot: Array<TokenSnapshot>;
  token: Token;
  totalPnl: Scalars['Decimal']['output'];
  updatedAt: Scalars['Time']['output'];
  valueUsd: Scalars['Decimal']['output'];
};

export type Token = {
  __typename?: 'Token';
  createdAt: Scalars['Time']['output'];
  currentUsdValue: Scalars['Decimal']['output'];
  historicalPrices: Array<HistoricalPrice>;
  id: Scalars['UUID']['output'];
  lastUpdated: Scalars['Time']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type TokenSnapshot = {
  __typename?: 'TokenSnapshot';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  valueUsd: Scalars['Decimal']['output'];
};

export type Transaction = {
  __typename?: 'Transaction';
  amount: Scalars['Decimal']['output'];
  category: TransactionCategory;
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  token: Token;
  transactionDate: Scalars['Time']['output'];
  transactionType: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  valueUsdAtTransaction: Scalars['Decimal']['output'];
};

export type TransactionCategory = {
  __typename?: 'TransactionCategory';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  transactions: Array<Transaction>;
  updatedAt: Scalars['Time']['output'];
};

export type User = {
  __typename?: 'User';
  alerts: Array<Alert>;
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  id: Scalars['UUID']['output'];
  metrics: Array<PortfolioMetric>;
  name: Scalars['String']['output'];
  roles: Array<Role>;
  updatedAt: Scalars['Time']['output'];
  wallets: Array<Wallet>;
};

export type Wallet = {
  __typename?: 'Wallet';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  subwallets: Array<Maybe<Subwallet>>;
  totalBalance: Scalars['Float']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type CreateSubwalletInput = {
  address: Scalars['String']['input'];
  chainId: Scalars['UUID']['input'];
  name: Scalars['String']['input'];
  walletId: Scalars['UUID']['input'];
};

export type CreateWalletInput = {
  name: Scalars['String']['input'];
};

export type WalletDetailQueryQueryVariables = Exact<{
  id: Scalars['UUID']['input'];
}>;


export type WalletDetailQueryQuery = { __typename?: 'Query', wallet: { __typename?: 'Wallet', createdAt: any, id: any, updatedAt: any, name: string, totalBalance: number, subwallets: Array<{ __typename?: 'Subwallet', id: any, createdAt: any, updatedAt: any, name: string, tokens: Array<{ __typename?: 'SubwalletToken', amount: any, valueUsd: any, totalPnl: any } | null>, snapshots: Array<{ __typename?: 'Snapshot', snapshotDate: any, totalPnl: any, totalValue: any, id: any, createdAt: any } | null> } | null> } };

export class TypedDocumentString<TResult, TVariables>
  extends String
  implements DocumentTypeDecoration<TResult, TVariables>
{
  __apiType?: DocumentTypeDecoration<TResult, TVariables>['__apiType'];

  constructor(private value: string, public __meta__?: Record<string, any> | undefined) {
    super(value);
  }

  toString(): string & DocumentTypeDecoration<TResult, TVariables> {
    return this.value;
  }
}

export const WalletDetailQueryDocument = new TypedDocumentString(`
    query WalletDetailQuery($id: UUID!) {
  wallet(id: $id) {
    subwallets {
      id
      createdAt
      updatedAt
      name
      tokens {
        amount
        valueUsd
        totalPnl
      }
      snapshots {
        snapshotDate
        totalPnl
        totalValue
        id
        createdAt
      }
    }
    createdAt
    id
    updatedAt
    name
    totalBalance
  }
}
    `) as unknown as TypedDocumentString<WalletDetailQueryQuery, WalletDetailQueryQueryVariables>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Decimal: { input: any; output: any; }
  Time: { input: any; output: any; }
  UUID: { input: any; output: any; }
};

export type Alert = {
  __typename?: 'Alert';
  condition: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  notificationSettings?: Maybe<Scalars['String']['output']>;
  token: Token;
  updatedAt: Scalars['Time']['output'];
  user: User;
};

export type AuthResponse = {
  __typename?: 'AuthResponse';
  token: Scalars['String']['output'];
  user: User;
};

export type Chain = {
  __typename?: 'Chain';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  subwallets: Array<Subwallet>;
  updatedAt: Scalars['Time']['output'];
};

export type CreateAlertInput = {
  condition: Scalars['String']['input'];
  notificationSettings?: InputMaybe<Scalars['String']['input']>;
  tokenId: Scalars['UUID']['input'];
};

export type CreateUserInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type HistoricalPrice = {
  __typename?: 'HistoricalPrice';
  createdAt: Scalars['Time']['output'];
  date: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  price: Scalars['Decimal']['output'];
  token: Token;
};

export type LoginInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  assignRole: User;
  createAlert: Alert;
  createSubwallet: Subwallet;
  createWallet: Wallet;
  deleteAlert: Scalars['Boolean']['output'];
  login: AuthResponse;
  register: AuthResponse;
  removeRole: User;
};


export type MutationAssignRoleArgs = {
  roleName: Scalars['String']['input'];
  userId: Scalars['UUID']['input'];
};


export type MutationCreateAlertArgs = {
  input: CreateAlertInput;
};


export type MutationCreateSubwalletArgs = {
  input: CreateSubwalletInput;
};


export type MutationCreateWalletArgs = {
  input: CreateWalletInput;
};


export type MutationDeleteAlertArgs = {
  id: Scalars['UUID']['input'];
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationRemoveRoleArgs = {
  roleName: Scalars['String']['input'];
  userId: Scalars['UUID']['input'];
};

export type PortfolioMetric = {
  __typename?: 'PortfolioMetric';
  calculatedAt: Scalars['Time']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  metricName: Scalars['String']['output'];
  metricValue: Scalars['Decimal']['output'];
  user: User;
};

export type PortfolioStats = {
  __typename?: 'PortfolioStats';
  dailyChange: Scalars['Decimal']['output'];
  topPerformers: Array<SubwalletToken>;
  totalPnl: Scalars['Decimal']['output'];
  totalValue: Scalars['Decimal']['output'];
  worstPerformers: Array<SubwalletToken>;
};

export type Query = {
  __typename?: 'Query';
  alerts: Array<Alert>;
  chain: Chain;
  chains: Array<Chain>;
  me: User;
  portfolioMetrics: Array<PortfolioMetric>;
  portfolioStats: PortfolioStats;
  roles: Array<Role>;
  subwallet: Subwallet;
  token: Token;
  tokenPriceHistory: Array<HistoricalPrice>;
  tokens: Array<Token>;
  transactionCategories: Array<TransactionCategory>;
  transactions: Array<Transaction>;
  wallet: Wallet;
  wallets: Array<Maybe<Wallet>>;
};


export type QueryChainArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryPortfolioMetricsArgs = {
  from: Scalars['Time']['input'];
  to: Scalars['Time']['input'];
};


export type QuerySubwalletArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryTokenArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryTokenPriceHistoryArgs = {
  from: Scalars['Time']['input'];
  id: Scalars['UUID']['input'];
  to: Scalars['Time']['input'];
};


export type QueryWalletArgs = {
  id: Scalars['UUID']['input'];
};

export type RegisterInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Role = {
  __typename?: 'Role';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  users: Array<User>;
};

export type Snapshot = {
  __typename?: 'Snapshot';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  snapshotDate: Scalars['Time']['output'];
  subwallet: Subwallet;
  totalPnl: Scalars['Decimal']['output'];
  totalValue: Scalars['Decimal']['output'];
};

export type Subwallet = {
  __typename?: 'Subwallet';
  address: Scalars['String']['output'];
  chain: Chain;
  createdAt: Scalars['Time']['output'];
  currentValue: Scalars['Float']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  snapshots: Array<Maybe<Snapshot>>;
  tokens: Array<Maybe<SubwalletToken>>;
  updatedAt: Scalars['Time']['output'];
};

export type SubwalletToken = {
  __typename?: 'SubwalletToken';
  amount: Scalars['Decimal']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  snapshot: Array<TokenSnapshot>;
  token: Token;
  totalPnl: Scalars['Decimal']['output'];
  updatedAt: Scalars['Time']['output'];
  valueUsd: Scalars['Decimal']['output'];
};

export type Token = {
  __typename?: 'Token';
  createdAt: Scalars['Time']['output'];
  currentUsdValue: Scalars['Decimal']['output'];
  historicalPrices: Array<HistoricalPrice>;
  id: Scalars['UUID']['output'];
  lastUpdated: Scalars['Time']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type TokenSnapshot = {
  __typename?: 'TokenSnapshot';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  valueUsd: Scalars['Decimal']['output'];
};

export type Transaction = {
  __typename?: 'Transaction';
  amount: Scalars['Decimal']['output'];
  category: TransactionCategory;
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  token: Token;
  transactionDate: Scalars['Time']['output'];
  transactionType: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  valueUsdAtTransaction: Scalars['Decimal']['output'];
};

export type TransactionCategory = {
  __typename?: 'TransactionCategory';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  transactions: Array<Transaction>;
  updatedAt: Scalars['Time']['output'];
};

export type User = {
  __typename?: 'User';
  alerts: Array<Alert>;
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  id: Scalars['UUID']['output'];
  metrics: Array<PortfolioMetric>;
  name: Scalars['String']['output'];
  roles: Array<Role>;
  updatedAt: Scalars['Time']['output'];
  wallets: Array<Wallet>;
};

export type Wallet = {
  __typename?: 'Wallet';
  createdAt: Scalars['Time']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  subwallets: Array<Maybe<Subwallet>>;
  totalBalance: Scalars['Float']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type CreateSubwalletInput = {
  address: Scalars['String']['input'];
  chainId: Scalars['UUID']['input'];
  name: Scalars['String']['input'];
  walletId: Scalars['UUID']['input'];
};

export type CreateWalletInput = {
  name: Scalars['String']['input'];
};

export type WalletDetailQueryQueryVariables = Exact<{
  id: Scalars['UUID']['input'];
}>;


export type WalletDetailQueryQuery = { __typename?: 'Query', wallet: { __typename?: 'Wallet', createdAt: any, id: any, updatedAt: any, name: string, totalBalance: number, subwallets: Array<{ __typename?: 'Subwallet', id: any, createdAt: any, updatedAt: any, name: string, tokens: Array<{ __typename?: 'SubwalletToken', amount: any, valueUsd: any, totalPnl: any } | null>, snapshots: Array<{ __typename?: 'Snapshot', snapshotDate: any, totalPnl: any, totalValue: any, id: any, createdAt: any } | null> } | null> } };
