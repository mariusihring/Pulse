import {InnerInstruction, Reward, TokenBalance, Transaction} from "@/lib/gql/graphql";

interface TokenTransferDetail {
    tokenMint: string;
    amount: number;
    direction: 'sent' | 'received';
    sender: string;
    recipient: string;
}

interface EnhancedTransaction {
    originalTransaction: Transaction;
    info: {
        type: 'transfer' | 'swap' | 'unknown';
        lamportTransfer: number;
        overallTransferDirection: 'sent' | 'received' | 'none' | 'mixed';
        transactionStatus: 'completed' | 'pending' | 'error';
        tokenTransfers: Array<{
            programId: string;
            direction: 'sent' | 'received' | 'unknown';
            instructionData: string;
            involvedAddresses: string[];
        }>;
        swapDetails: Array<{
            programId: string;
            instructionData: string;
            involvedAddresses: string[];
            swapper: string;
        }>;
        fee: number;
        additionalLogs: string[];
        innerInstructions: InnerInstruction[];
        tokenPreBalances: TokenBalance[];
        tokenPostBalances: TokenBalance[];
        rewards: Reward[];
        tokenTransferDetails: TokenTransferDetail[];
    };
}

function determineOverallTransferDirection(
    lamportDelta: number,
    tokenDetails: TokenTransferDetail[]
): 'sent' | 'received' | 'none' | 'mixed' {
    if (lamportDelta > 0) return 'received';
    if (lamportDelta < 0) return 'sent';
    if (tokenDetails.length === 0) return 'none';
    const directions = new Set(tokenDetails.map(td => td.direction));
    return directions.size === 1 ? directions.values().next().value : 'mixed';
}

function determineTransactionStatus(result: {
    block_time?: string | null;
    meta: { status?: { error_message?: string; ok?: string } };
}): 'completed' | 'pending' | 'error' {
    if (result.meta.status && result.meta.status.error_message) {
        return 'error';
    }
    return result.block_time ? 'completed' : 'pending';
}

export function parseTransaction(transaction: Transaction, walletAddress: string): EnhancedTransaction {
    if (!transaction.result) {
        return {
            originalTransaction: transaction,
            info: {
                type: 'unknown',
                lamportTransfer: 0,
                overallTransferDirection: 'none',
                transactionStatus: 'pending',
                tokenTransfers: [],
                swapDetails: [],
                fee: 0,
                additionalLogs: [],
                innerInstructions: [],
                tokenPreBalances: [],
                tokenPostBalances: [],
                rewards: [],
                tokenTransferDetails: [],
            },
        };
    }

    const { meta, transaction: txData } = transaction.result;
    const message = txData.message;
    const accountKeys = message.account_keys;

    // Calculate lamport delta for the wallet.
    const walletIndex = accountKeys.indexOf(walletAddress);
    const lamportTransfer = walletIndex >= 0
        ? meta.post_balances[walletIndex] - meta.pre_balances[walletIndex]
        : 0;

    // Define known program IDs.
    const SYSTEM_PROGRAM_ID = "11111111111111111111111111111111";
    const TOKEN_PROGRAM_ID = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA";
    const SWAP_PROGRAM_IDS = new Set(["JUP4Fb6vV9N6Pz...", "RadyumSwapProgramId"]);
    const SWAP_PROGRAM_NAMES: { [key: string]: string } = {
        "JUP4Fb6vV9N6Pz...": "Jupiter",
        "RadyumSwapProgramId": "Radyum",
    };

    let detectedType: 'transfer' | 'swap' | 'unknown' = 'unknown';
    const tokenTransfers: EnhancedTransaction['info']['tokenTransfers'] = [];
    const swapDetails: EnhancedTransaction['info']['swapDetails'] = [];

    // Process instructions in one pass.
    for (const instr of message.instructions) {
        const programId = accountKeys[instr.program_id_index];
        const involvedAddresses = instr.accounts.map(idx => accountKeys[idx]);

        if (SWAP_PROGRAM_IDS.has(programId)) {
            detectedType = 'swap';
            swapDetails.push({
                programId,
                instructionData: instr.data,
                involvedAddresses,
                swapper: SWAP_PROGRAM_NAMES[programId] || programId,
            });
        } else if (programId === SYSTEM_PROGRAM_ID || programId === TOKEN_PROGRAM_ID) {
            if (involvedAddresses.includes(walletAddress)) {
                if (detectedType === 'unknown') detectedType = 'transfer';
                const direction = involvedAddresses[0] === walletAddress ? 'sent' : 'received';
                tokenTransfers.push({
                    programId,
                    direction,
                    instructionData: instr.data,
                    involvedAddresses,
                });
            }
        }
    }

    // Fallback using log messages if type remains unknown.
    if (detectedType === 'unknown' && meta.log_messages?.length) {
        for (const log of meta.log_messages) {
            const lowerLog = log.toLowerCase();
            if (lowerLog.includes('swap')) { detectedType = 'swap'; break; }
            if (lowerLog.includes('transfer')) { detectedType = 'transfer'; break; }
        }
    }

    // Compute token balance deltas using a Map.
    const tokenDeltaMap = new Map<string, number>();
    (meta.pre_token_balances || []).forEach(tb => {
        if (tb.owner === walletAddress) {
            tokenDeltaMap.set(
                tb.mint,
                (tokenDeltaMap.get(tb.mint) || 0) - tb.ui_token_amount.ui_amount
            );
        }
    });
    (meta.post_token_balances || []).forEach(tb => {
        if (tb.owner === walletAddress) {
            tokenDeltaMap.set(
                tb.mint,
                (tokenDeltaMap.get(tb.mint) || 0) + tb.ui_token_amount.ui_amount
            );
        }
    });
    const tokenTransferDetails: TokenTransferDetail[] = Array.from(tokenDeltaMap.entries())
        .filter(([, delta]) => delta !== 0)
        .map(([mint, delta]) => ({
            tokenMint: mint,
            amount: Math.abs(delta),
            direction: delta > 0 ? 'received' : 'sent',
            sender: delta < 0 ? walletAddress : "unknown",
            recipient: delta > 0 ? walletAddress : "unknown",
        }));

    const overallTransferDirection = determineOverallTransferDirection(lamportTransfer, tokenTransferDetails);
    const transactionStatus = determineTransactionStatus(transaction.result);

    const enhancedInfo = {
        type: detectedType,
        lamportTransfer,
        overallTransferDirection,
        transactionStatus,
        tokenTransfers,
        swapDetails,
        fee: meta.fee,
        additionalLogs: meta.log_messages || [],
        innerInstructions: meta.inner_instructions || [],
        tokenPreBalances: meta.pre_token_balances || [],
        tokenPostBalances: meta.post_token_balances || [],
        rewards: meta.rewards || [],
        tokenTransferDetails,
    };

    return {
        originalTransaction: transaction,
        info: enhancedInfo,
    };
}
