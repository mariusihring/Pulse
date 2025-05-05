<?php

namespace App\Services;

use App\Models\Token;
use App\Models\TokenHoldings;
use App\Models\TokenSwap;
use App\Models\Wallet;
use App\Models\WalletSnapshot;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Symfony\Component\Uid\Uuid;

// Maybe check here for more good apis
// https://github.com/public-apis/public-apis?tab=readme-ov-file#cryptocurrency

class WalletService
{
    private Client $client;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function loadPortfolio(string $userId, string $walletAddress, string $chainId): array
    {
        $portfolioData = $this->fetchPortfolioData($walletAddress);
        $wallet = $this->createOrUpdateWallet($userId, $walletAddress, $chainId, $portfolioData);
        
        $mintPrices = $this->fetchTokenPrices($portfolioData['tokens']);
        $walletValue = $this->processTokens($portfolioData['tokens'], $mintPrices, $wallet, $userId, $chainId);
        
        $this->updateWalletValue($wallet, $walletValue);
        $this->createWalletSnapshot($wallet, $walletValue);

        return $portfolioData;
    }

    private function fetchPortfolioData(string $walletAddress): array
    {
        $response = $this->client->get(sprintf("/account/mainnet/%s/portfolio", $walletAddress));
        return json_decode($response->getBody(), true);
    }

private function createOrUpdateWallet(string $userId, string $walletAddress, string $chainId, array $portfolioData): Wallet
{
    return Wallet::updateOrCreate(
        ['address' => $walletAddress],
        [
            'chain_token_amount' => $portfolioData['nativeBalance']['solana'],
            'value' => 0,
            'chain_id' => $chainId,
            'user_id' => $userId,
        ],
        [
            'name' => 'My New Wallet',
            'favorite' => false,
        ]
    );
}

    private function fetchTokenPrices(array $tokens): array
    {
        $mints = array_column($tokens, 'mint');
        $response = $this->client->post('/token/mainnet/prices', [
            'json' => ['addresses' => $mints],
        ]);

        $prices = json_decode($response->getBody(), true);
        return array_column($prices, null, 'tokenAddress');
    }

    private function processTokens(array $tokens, array $mintPrices, Wallet $wallet, string $userId, string $chainId): float
    {
        $walletValue = 0;

        foreach ($tokens as $token) {
            $tokenModel = $this->createOrUpdateToken($token, $mintPrices[$token['mint']], $chainId);
            $tokenValue = $this->createTokenHolding($token, $tokenModel, $wallet, $userId);
            $walletValue += $tokenValue;
        }

        return $walletValue;
    }

    private function createOrUpdateToken(array $tokenData, array $priceData, string $chainId): Token
    {
        try {
            $token = Token::where('mint', $tokenData['associatedTokenAddress'])->first();

            if ($token) {
                $token->current_price = $priceData['usdPrice'];
                $token->save();

                return $token;
            } else {
                $token = Token::create([
                    'name' => $tokenData['name'],
                    'chain_id' => $chainId,
                    'current_price' => $priceData['usdPrice'],
                    'address' => $tokenData['associatedTokenAddress'],
                    'symbol' => $tokenData['symbol'],
                    'logo' => $tokenData['logo'] ?? null,
                    'mint' => $tokenData['associatedTokenAddress'],
                ]);

                return $token;
            }
        } catch (\Exception $e) {
            \Log::error("Error creating or updating token: " . $e->getMessage());
            throw $e;
        }
    }

    private function createTokenHolding(array $token, Token $tokenModel, Wallet $wallet, string $userId): float
    {
        $value = $token['amount'] * $tokenModel->current_price;
        
        TokenHoldings::updateOrCreate(
            [
                'user_id' => $userId,
                'token_id' => $tokenModel->id,
                'wallet_id' => $wallet->id,
            ],
            [
                'amount' => $token['amount'],
                'value' => $value,
            ]
        );

        return $value;
    }

    private function updateWalletValue(Wallet $wallet, float $value): void
    {
        $wallet->update(['value' => $value]);
    }

    private function createWalletSnapshot(Wallet $wallet, float $value): void
    {
        WalletSnapshot::create([
            'value' => $value,
            'wallet_id' => $wallet->id,
        ]);
    }

    public function refreshWallet(string $walletAddress): Wallet
    {
        $wallet = Wallet::where('address', $walletAddress)->firstOrFail();
        $wallet->refresh();
        return $wallet;
    }

    public function getTokenSwaps(string $walletAddress, string $chainId): array
    {
        $wallet = Wallet::where('address', $walletAddress)->firstOrFail();
        $swaps = $this->fetchTokenSwaps($walletAddress);
        
        if (empty($swaps['result'])) {
            return [];
        }

        return $this->processTokenSwaps($swaps['result'], $wallet, $chainId);
    }

    private function fetchTokenSwaps(string $walletAddress): array
    {
        $response = $this->client->get(sprintf("/account/mainnet/%s/swaps?order=DESC", $walletAddress));
        return json_decode($response->getBody(), true);
    }

    private function processTokenSwaps(array $swaps, Wallet $wallet, string $chainId): array
    {
        $createdSwaps = [];

        foreach ($swaps as $swapData) {
            $token = $this->getOrCreateTokenForSwap($swapData, $chainId);
            $createdSwaps[] = $this->createTokenSwap($swapData, $token, $wallet, $chainId);
        }

        return $createdSwaps;
    }

    private function getOrCreateTokenForSwap(array $swapData, string $chainId): Token
    {
        $tokenData = $swapData['bought']['address'] === $swapData['baseToken'] 
            ? $swapData['bought'] 
            : $swapData['sold'];

        return Token::firstOrCreate(
            ['address' => $swapData['baseToken']],
            [
                'chain_id' => $chainId,
                'name' => $tokenData['name'],
                'current_price' => $tokenData['usdPrice'],
                'logo' => $tokenData['logo'] ?? null,
                'symbol' => $tokenData['symbol'],
                'mint' => $tokenData['address'],
            ]
        );
    }

    private function createTokenSwap(array $swapData, Token $token, Wallet $wallet, string $chainId): TokenSwap
    {
        return TokenSwap::updateOrCreate(
            ['transaction_hash' => $swapData['transactionHash']],
            [
                'chain_id' => $chainId,
                'token_id' => $token->id,
                'wallet_id' => $wallet->id,
                'transaction_hash' => $swapData['transactionHash'],
                'transaction_type' => $swapData['transactionType'],
                'transaction_index' => $swapData['transactionIndex'],
                'sub_category' => $swapData['subCategory'] ?? null,
                'block_timestamp' => $swapData['blockTimestamp'],
                'block_number' => $swapData['blockNumber'] ?: 0,
                'wallet_address' => $swapData['walletAddress'],
                'pair_address' => $swapData['pairAddress'],
                'pair_label' => $swapData['pairLabel'],
                'exchange_address' => $swapData['exchangeAddress'],
                'exchange_name' => $swapData['exchangeName'],
                'exchange_logo' => $swapData['exchangeLogo'] ?? null,
                'base_token' => $swapData['baseToken'],
                'quote_token' => $swapData['quoteToken'],
                'bought' => $this->formatTokenData($swapData['bought']),
                'sold' => $this->formatTokenData($swapData['sold']),
                'base_quote_price' => $swapData['baseQuotePrice'],
                'total_value_usd' => $swapData['totalValueUsd'],
            ]
        );
    }

    private function formatTokenData(array $token): array
    {
        return [
            'address' => $token['address'],
            'amount' => $token['amount'],
            'usdPrice' => $token['usdPrice'],
            'usdAmount' => $token['usdAmount'],
            'symbol' => $token['symbol'],
            'logo' => $token['logo'] ?? null,
            'name' => $token['name'],
            'tokenType' => $token['tokenType'],
        ];
    }

    public function getTransactions(string $walletAddress)
    {
        //TODO: we get transactions and if they have tokenswap as type we add a foreign key to the transaction so we can load the swap
        todo("do this");
    }
}
