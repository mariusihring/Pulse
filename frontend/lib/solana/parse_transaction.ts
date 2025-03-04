import {TransactionResult} from "@/lib/gql/graphql";

type ParsedTransaction = {
    type: 'send' | 'receive' | 'swap' | 'unknown';
    from: string;
    to: string;
    amount: number; // in lamports or token amount (if applicable)
    fee: number;
    details: {
        preBalance: number;
        postBalance: number;
        tokenChanges?: TokenChange[];
        instructions: ParsedInstruction[];
    };
};

type TokenChange = {
    accountIndex: number;
    mint: string;
    preAmount: string;
    postAmount: string;
    change: number;
};

const SYSTEM_PROGRAM_ID = '11111111111111111111111111111111';
const COMPUTE_BUDGET_PROGRAM_ID = 'ComputeBudget111111111111111111111111111111';

// Update this list with known swap program IDs (example placeholders)
const KNOWN_SWAP_PROGRAM_IDS = [
    'Swap111111111111111111111111111111111111',
    'AnotherSwap11111111111111111111111111111'
];

/**
 * Parses a transaction for a given address.
 * It examines lamport balance changes, token balance changes, and instructions
 * to determine if the transaction was a simple send/receive or a swap.
 *
 * @param address - The public key string for the account of interest.
 * @param tx - The parsed TransactionResponse from Solana RPC.
 * @returns A ParsedTransaction with details.
 */
function parseTransactionForAddress(
    address: string,
    tx: TransactionResult
): ParsedTransaction {
    const accountKeys: string[] = tx.transaction.message.account_keys;
    const meta = tx.meta;
    if (!meta) {
        throw new Error('Missing metadata in transaction');
    }
    if (!meta.pre_balances || !meta.post_balances) {
        throw new Error('Missing lamport balance information');
    }

    // Find the index of the given address in account_keys
    const accountIndex = accountKeys.findIndex((key) => key === address);
    if (accountIndex === -1) {
        throw new Error('The provided address is not involved in this transaction');
    }

    // Check lamport changes for the given address
    const preLamport = meta.pre_balances[accountIndex];
    const postLamport = meta.post_balances[accountIndex];
    const lamportChange = postLamport - preLamport;

    // Create a holder for token balance changes (if any)
    const tokenChanges: TokenChange[] = [];
    if (meta.pre_token_balances && meta.post_token_balances &&
        meta.pre_token_balances.length > 0) {
        // Both pre and post token balance arrays are available
        // They usually include entries for all token accounts touched in the transaction
        for (let i = 0; i < meta.pre_token_balances.length; i++) {
            const preToken = meta.pre_token_balances[i];
            const postToken = meta.post_token_balances[i];
            // Ensure both entries refer to the same token account
            if (preToken.account_index === postToken.account_index) {
                const change =
                    parseFloat(postToken.ui_token_amount.amount) -
                    parseFloat(preToken.ui_token_amount.amount);
                tokenChanges.push({
                    accountIndex: preToken.account_index,
                    mint: preToken.mint,
                    preAmount: preToken.ui_token_amount.amount,
                    postAmount: postToken.ui_token_amount.amount,
                    change,
                });
            }
        }
    }

    // Start by assuming a basic type based on lamport change
    let type: 'send' | 'receive' | 'swap' | 'unknown' = 'unknown';
    let from = '';
    let to = '';

    if (lamportChange < 0) {
        type = 'send';
        from = address;
        // Look for the first account with a positive lamport change as a candidate recipient
        const possibleRecipients = accountKeys.filter((_, idx) => {
            return meta.post_balances[idx] - meta.pre_balances[idx] > 0;
        });
        to = possibleRecipients.length > 0 ? possibleRecipients[0] : '';
    } else if (lamportChange > 0) {
        type = 'receive';
        to = address;
        // Look for the first account with a negative lamport change as a candidate sender
        const possibleSenders = accountKeys.filter((_, idx) => {
            return meta.post_balances[idx] - meta.pre_balances[idx] < 0;
        });
        from = possibleSenders.length > 0 ? possibleSenders[0] : '';
    }

    // Check instructions for known swap program calls.
    // We also check if token changes exist that indicate a two-sided token movement.
    const instructions = tx.transaction.message.instructions;
    const invokedProgramIds = instructions.map(instr => accountKeys[instr.program_id_index]);

    const hasSwapProgram = invokedProgramIds.some((programId) =>
        KNOWN_SWAP_PROGRAM_IDS.includes(programId)
    );

    // Alternatively, if there are at least two token changes (one positive and one negative), flag as a swap.
    const tokenSend = tokenChanges.find(tc => tc.change < 0);
    const tokenReceive = tokenChanges.find(tc => tc.change > 0);
    const isTokenSwap = tokenChanges.length >= 2 && !!tokenSend && !!tokenReceive;

    if (hasSwapProgram || isTokenSwap) {
        type = 'swap';
        // For swaps, set the sending address as the one with negative token change,
        // and receiving as the one with positive token change.
        // Note: in many cases the same wallet may be involved on both sides.
        from = tokenSend ? accountKeys[tokenSend.accountIndex] : from;
        to = tokenReceive ? accountKeys[tokenReceive.accountIndex] : to;
    }

    // Return a neat formatted result.
    return {
        type,
        from,
        to,
        amount: Math.abs(lamportChange),
        fee: meta.fee,
        details: {
            preBalance: preLamport,
            postBalance: postLamport,
            tokenChanges: tokenChanges.length ? tokenChanges : undefined,
            instructions,
        },
    };
}

