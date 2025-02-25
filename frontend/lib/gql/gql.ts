/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

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
    "\n    mutation StartWalletScan($input: String!) {\n        startWalletUpdate(walletAddress: $input) {\n            id\n            walletAddress\n        }\n    }\n": typeof types.StartWalletScanDocument,
    "\n    subscription WalletUpdate($input: ID!) {\n        walletUpdates(jobID: $input) {\n            JobID\n            Progress\n            Wallet {\n                address\n                sol_balance\n                sol_value\n                wallet_value\n                last_updated\n                tokens {\n                    name\n                    address\n                    pool\n                    description\n                    image\n                    amount\n                    price\n                    pnl\n                    invested\n                    value\n                    history_prices\n                }\n                transactions {\n                    jsonrpc\n                    id\n                    result {\n                        block_time\n                        slot\n                        meta {\n                            compute_units_consumed\n                            fee\n                            log_messages\n                            post_balances\n                            pre_balances\n                            post_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            inner_instructions {\n                                index\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                            pre_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            rewards {\n                                info\n                            }\n                            status {\n                                ok\n                                error_message\n                            }\n                        }\n                        transaction {\n                            signatures\n                            message {\n                                account_keys\n                                recent_blockhash\n                                address_table_lookups {\n                                    account_key\n                                    readonly_indexes\n                                    writable_indexes\n                                }\n                                header {\n                                    num_readonly_signed_accounts\n                                    num_readonly_unsigned_accounts\n                                    num_required_signatures\n                                }\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                        }\n                    }\n                    err {\n                        code\n                        message\n                    }\n                }\n            }\n        }\n    }\n": typeof types.WalletUpdateDocument,
};
const documents: Documents = {
    "\n    mutation StartWalletScan($input: String!) {\n        startWalletUpdate(walletAddress: $input) {\n            id\n            walletAddress\n        }\n    }\n": types.StartWalletScanDocument,
    "\n    subscription WalletUpdate($input: ID!) {\n        walletUpdates(jobID: $input) {\n            JobID\n            Progress\n            Wallet {\n                address\n                sol_balance\n                sol_value\n                wallet_value\n                last_updated\n                tokens {\n                    name\n                    address\n                    pool\n                    description\n                    image\n                    amount\n                    price\n                    pnl\n                    invested\n                    value\n                    history_prices\n                }\n                transactions {\n                    jsonrpc\n                    id\n                    result {\n                        block_time\n                        slot\n                        meta {\n                            compute_units_consumed\n                            fee\n                            log_messages\n                            post_balances\n                            pre_balances\n                            post_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            inner_instructions {\n                                index\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                            pre_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            rewards {\n                                info\n                            }\n                            status {\n                                ok\n                                error_message\n                            }\n                        }\n                        transaction {\n                            signatures\n                            message {\n                                account_keys\n                                recent_blockhash\n                                address_table_lookups {\n                                    account_key\n                                    readonly_indexes\n                                    writable_indexes\n                                }\n                                header {\n                                    num_readonly_signed_accounts\n                                    num_readonly_unsigned_accounts\n                                    num_required_signatures\n                                }\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                        }\n                    }\n                    err {\n                        code\n                        message\n                    }\n                }\n            }\n        }\n    }\n": types.WalletUpdateDocument,
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
export function graphql(source: "\n    mutation StartWalletScan($input: String!) {\n        startWalletUpdate(walletAddress: $input) {\n            id\n            walletAddress\n        }\n    }\n"): (typeof documents)["\n    mutation StartWalletScan($input: String!) {\n        startWalletUpdate(walletAddress: $input) {\n            id\n            walletAddress\n        }\n    }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    subscription WalletUpdate($input: ID!) {\n        walletUpdates(jobID: $input) {\n            JobID\n            Progress\n            Wallet {\n                address\n                sol_balance\n                sol_value\n                wallet_value\n                last_updated\n                tokens {\n                    name\n                    address\n                    pool\n                    description\n                    image\n                    amount\n                    price\n                    pnl\n                    invested\n                    value\n                    history_prices\n                }\n                transactions {\n                    jsonrpc\n                    id\n                    result {\n                        block_time\n                        slot\n                        meta {\n                            compute_units_consumed\n                            fee\n                            log_messages\n                            post_balances\n                            pre_balances\n                            post_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            inner_instructions {\n                                index\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                            pre_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            rewards {\n                                info\n                            }\n                            status {\n                                ok\n                                error_message\n                            }\n                        }\n                        transaction {\n                            signatures\n                            message {\n                                account_keys\n                                recent_blockhash\n                                address_table_lookups {\n                                    account_key\n                                    readonly_indexes\n                                    writable_indexes\n                                }\n                                header {\n                                    num_readonly_signed_accounts\n                                    num_readonly_unsigned_accounts\n                                    num_required_signatures\n                                }\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                        }\n                    }\n                    err {\n                        code\n                        message\n                    }\n                }\n            }\n        }\n    }\n"): (typeof documents)["\n    subscription WalletUpdate($input: ID!) {\n        walletUpdates(jobID: $input) {\n            JobID\n            Progress\n            Wallet {\n                address\n                sol_balance\n                sol_value\n                wallet_value\n                last_updated\n                tokens {\n                    name\n                    address\n                    pool\n                    description\n                    image\n                    amount\n                    price\n                    pnl\n                    invested\n                    value\n                    history_prices\n                }\n                transactions {\n                    jsonrpc\n                    id\n                    result {\n                        block_time\n                        slot\n                        meta {\n                            compute_units_consumed\n                            fee\n                            log_messages\n                            post_balances\n                            pre_balances\n                            post_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            inner_instructions {\n                                index\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                            pre_token_balances {\n                                account_index\n                                mint\n                                owner\n                                program_id\n                                ui_token_amount {\n                                    amount\n                                    decimals\n                                    ui_amount\n                                    ui_amount_string\n                                }\n                            }\n                            rewards {\n                                info\n                            }\n                            status {\n                                ok\n                                error_message\n                            }\n                        }\n                        transaction {\n                            signatures\n                            message {\n                                account_keys\n                                recent_blockhash\n                                address_table_lookups {\n                                    account_key\n                                    readonly_indexes\n                                    writable_indexes\n                                }\n                                header {\n                                    num_readonly_signed_accounts\n                                    num_readonly_unsigned_accounts\n                                    num_required_signatures\n                                }\n                                instructions {\n                                    accounts\n                                    data\n                                    program_id_index\n                                    stack_height\n                                }\n                            }\n                        }\n                    }\n                    err {\n                        code\n                        message\n                    }\n                }\n            }\n        }\n    }\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;