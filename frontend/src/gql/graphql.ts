/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
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

export type AddressTableLookup = {
  __typename?: 'AddressTableLookup';
  account_key: Scalars['String']['output'];
  readonly_indexes: Array<Scalars['Int']['output']>;
  writable_indexes: Array<Scalars['Int']['output']>;
};

export type AuthResponse = {
  __typename?: 'AuthResponse';
  token: Scalars['String']['output'];
  user: User;
};

export type CreateUserInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Error = {
  __typename?: 'Error';
  code: Scalars['Int']['output'];
  message: Scalars['String']['output'];
};

export type InnerInstruction = {
  __typename?: 'InnerInstruction';
  index: Scalars['Int']['output'];
  instructions: Array<Instruction>;
};

export type Instruction = {
  __typename?: 'Instruction';
  accounts: Array<Scalars['Int']['output']>;
  data: Scalars['String']['output'];
  program_id_index: Scalars['Int']['output'];
  stack_height: Scalars['Int']['output'];
};

export type Job = {
  __typename?: 'Job';
  id: Scalars['ID']['output'];
  walletAddress: Scalars['String']['output'];
};

export type LoginInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type MessageHeader = {
  __typename?: 'MessageHeader';
  num_readonly_signed_accounts: Scalars['Int']['output'];
  num_readonly_unsigned_accounts: Scalars['Int']['output'];
  num_required_signatures: Scalars['Int']['output'];
};

export type Meta = {
  __typename?: 'Meta';
  compute_units_consumed: Scalars['Int']['output'];
  fee: Scalars['Int']['output'];
  inner_instructions: Array<InnerInstruction>;
  log_messages: Array<Scalars['String']['output']>;
  post_balances: Array<Scalars['Int']['output']>;
  post_token_balances: Array<TokenBalance>;
  pre_balances: Array<Scalars['Int']['output']>;
  pre_token_balances: Array<TokenBalance>;
  rewards: Array<Reward>;
  status: Status;
};

export type Mutation = {
  __typename?: 'Mutation';
  login: AuthResponse;
  register: AuthResponse;
  startWalletUpdate: Job;
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationStartWalletUpdateArgs = {
  walletAddress: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  me: User;
};

export type RegisterInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Reward = {
  __typename?: 'Reward';
  info: Scalars['String']['output'];
};

export type Status = {
  __typename?: 'Status';
  error_message?: Maybe<Scalars['String']['output']>;
  ok?: Maybe<Scalars['String']['output']>;
};

export type StatusMessage = {
  __typename?: 'StatusMessage';
  status: Scalars['String']['output'];
};

export type Subscription = {
  __typename?: 'Subscription';
  walletUpdates: WalletUpdate;
};


export type SubscriptionWalletUpdatesArgs = {
  jobID: Scalars['ID']['input'];
};

export type Token = {
  __typename?: 'Token';
  address: Scalars['String']['output'];
  amount: Scalars['Float']['output'];
  description: Scalars['String']['output'];
  history_prices: Array<Scalars['Float']['output']>;
  image: Scalars['String']['output'];
  invested: Scalars['Float']['output'];
  name: Scalars['String']['output'];
  pnl: Scalars['Float']['output'];
  pool: Scalars['String']['output'];
  price: Scalars['Float']['output'];
  value: Scalars['Float']['output'];
};

export type TokenAmount = {
  __typename?: 'TokenAmount';
  amount: Scalars['String']['output'];
  decimals: Scalars['Int']['output'];
  ui_amount: Scalars['Float']['output'];
  ui_amount_string: Scalars['String']['output'];
};

export type TokenBalance = {
  __typename?: 'TokenBalance';
  account_index: Scalars['Int']['output'];
  mint: Scalars['String']['output'];
  owner: Scalars['String']['output'];
  program_id: Scalars['String']['output'];
  ui_token_amount: TokenAmount;
};

export type Transaction = {
  __typename?: 'Transaction';
  err?: Maybe<Error>;
  id: Scalars['Int']['output'];
  jsonrpc: Scalars['String']['output'];
  result?: Maybe<TransactionResult>;
};

export type TransactionData = {
  __typename?: 'TransactionData';
  message: TransactionMessage;
  signatures: Array<Scalars['String']['output']>;
};

export type TransactionMessage = {
  __typename?: 'TransactionMessage';
  account_keys: Array<Scalars['String']['output']>;
  address_table_lookups: Array<AddressTableLookup>;
  header: MessageHeader;
  instructions: Array<Instruction>;
  recent_blockhash: Scalars['String']['output'];
};

export type TransactionResult = {
  __typename?: 'TransactionResult';
  block_time: Scalars['Int']['output'];
  meta: Meta;
  slot: Scalars['Int']['output'];
  transaction: TransactionData;
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type Wallet = {
  __typename?: 'Wallet';
  address: Scalars['String']['output'];
  last_updated: Scalars['String']['output'];
  sol_balance: Scalars['Float']['output'];
  sol_value: Scalars['Float']['output'];
  tokens: Array<Token>;
  transactions: Array<Transaction>;
  wallet_value: Scalars['Float']['output'];
};

export type WalletUpdate = {
  __typename?: 'WalletUpdate';
  JobID: Scalars['ID']['output'];
  Progress: Scalars['Int']['output'];
  Wallet: Wallet;
};

export type StartWalletUpdateMutationVariables = Exact<{
  walletAddress: Scalars['String']['input'];
}>;


export type StartWalletUpdateMutation = { __typename?: 'Mutation', startWalletUpdate: { __typename?: 'Job', id: string, walletAddress: string } };

export type WalletUpdatesSubscriptionVariables = Exact<{
  jobID: Scalars['ID']['input'];
}>;


export type WalletUpdatesSubscription = { __typename?: 'Subscription', walletUpdates: { __typename?: 'WalletUpdate', JobID: string, Progress: number, Wallet: { __typename?: 'Wallet', address: string, sol_balance: number, sol_value: number, wallet_value: number, last_updated: string, tokens: Array<{ __typename?: 'Token', name: string, address: string, pool: string, description: string, image: string, amount: number, price: number, pnl: number, invested: number, value: number, history_prices: Array<number> }>, transactions: Array<{ __typename?: 'Transaction', jsonrpc: string, id: number }> } } };


export const StartWalletUpdateDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"StartWalletUpdate"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"walletAddress"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"startWalletUpdate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"walletAddress"},"value":{"kind":"Variable","name":{"kind":"Name","value":"walletAddress"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"walletAddress"}}]}}]}}]} as unknown as DocumentNode<StartWalletUpdateMutation, StartWalletUpdateMutationVariables>;
export const WalletUpdatesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"WalletUpdates"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"jobID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"walletUpdates"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"jobID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"jobID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"JobID"}},{"kind":"Field","name":{"kind":"Name","value":"Progress"}},{"kind":"Field","name":{"kind":"Name","value":"Wallet"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"sol_balance"}},{"kind":"Field","name":{"kind":"Name","value":"sol_value"}},{"kind":"Field","name":{"kind":"Name","value":"wallet_value"}},{"kind":"Field","name":{"kind":"Name","value":"last_updated"}},{"kind":"Field","name":{"kind":"Name","value":"tokens"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"pool"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"image"}},{"kind":"Field","name":{"kind":"Name","value":"amount"}},{"kind":"Field","name":{"kind":"Name","value":"price"}},{"kind":"Field","name":{"kind":"Name","value":"pnl"}},{"kind":"Field","name":{"kind":"Name","value":"invested"}},{"kind":"Field","name":{"kind":"Name","value":"value"}},{"kind":"Field","name":{"kind":"Name","value":"history_prices"}}]}},{"kind":"Field","name":{"kind":"Name","value":"transactions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"jsonrpc"}},{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]}}]}}]} as unknown as DocumentNode<WalletUpdatesSubscription, WalletUpdatesSubscriptionVariables>;