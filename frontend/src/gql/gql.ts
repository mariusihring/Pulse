/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n  mutation StartWalletUpdate($walletAddress: String!) {\n    startWalletUpdate(walletAddress: $walletAddress) {\n      id\n      walletAddress\n    }\n  }\n": typeof types.StartWalletUpdateDocument,
    "\n  subscription WalletUpdates ($jobID: ID!) {\n    walletUpdates(jobID: $jobID) {\n      JobID\n        Progress\n        Wallet {\n            address\n            sol_balance\n            sol_value\n            wallet_value\n            last_updated\n            tokens {\n                name\n                address\n                pool\n                description\n                image\n                amount\n                price\n                pnl\n                invested\n                value\n                history_prices\n            }\n            transactions {\n                jsonrpc\n                id\n            }\n        }\n    }\n}\n": typeof types.WalletUpdatesDocument,
};
const documents: Documents = {
    "\n  mutation StartWalletUpdate($walletAddress: String!) {\n    startWalletUpdate(walletAddress: $walletAddress) {\n      id\n      walletAddress\n    }\n  }\n": types.StartWalletUpdateDocument,
    "\n  subscription WalletUpdates ($jobID: ID!) {\n    walletUpdates(jobID: $jobID) {\n      JobID\n        Progress\n        Wallet {\n            address\n            sol_balance\n            sol_value\n            wallet_value\n            last_updated\n            tokens {\n                name\n                address\n                pool\n                description\n                image\n                amount\n                price\n                pnl\n                invested\n                value\n                history_prices\n            }\n            transactions {\n                jsonrpc\n                id\n            }\n        }\n    }\n}\n": types.WalletUpdatesDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation StartWalletUpdate($walletAddress: String!) {\n    startWalletUpdate(walletAddress: $walletAddress) {\n      id\n      walletAddress\n    }\n  }\n"): (typeof documents)["\n  mutation StartWalletUpdate($walletAddress: String!) {\n    startWalletUpdate(walletAddress: $walletAddress) {\n      id\n      walletAddress\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  subscription WalletUpdates ($jobID: ID!) {\n    walletUpdates(jobID: $jobID) {\n      JobID\n        Progress\n        Wallet {\n            address\n            sol_balance\n            sol_value\n            wallet_value\n            last_updated\n            tokens {\n                name\n                address\n                pool\n                description\n                image\n                amount\n                price\n                pnl\n                invested\n                value\n                history_prices\n            }\n            transactions {\n                jsonrpc\n                id\n            }\n        }\n    }\n}\n"): (typeof documents)["\n  subscription WalletUpdates ($jobID: ID!) {\n    walletUpdates(jobID: $jobID) {\n      JobID\n        Progress\n        Wallet {\n            address\n            sol_balance\n            sol_value\n            wallet_value\n            last_updated\n            tokens {\n                name\n                address\n                pool\n                description\n                image\n                amount\n                price\n                pnl\n                invested\n                value\n                history_prices\n            }\n            transactions {\n                jsonrpc\n                id\n            }\n        }\n    }\n}\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;