<?php

namespace App\Services;

use App\Models\Token;
use App\Models\TokenHoldings;
use App\Models\TokenSwap;
use App\Models\Wallet;
use App\Models\WalletSnapshot;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\DB;

// Maybe check here for more good apis
// https://github.com/public-apis/public-apis?tab=readme-ov-file#cryptocurrency

class WalletService
{
    public function loadPortfolio(string $userId, string $walletAddress, string $chainId)
    {
        $client = app(Client::class);
        $walletValue = 0;
        $response = $client->get(sprintf("/account/mainnet/%s/portfolio", $walletAddress));
        $body =  json_decode($response->getBody(), true);

        $wallet = new Wallet();
        $wallet->address = $walletAddress;
        $wallet->name = 'My New Wallet';
        $wallet->chain_token_amount = $body['nativeBalance']['solana'];
        $wallet->value = 0;
        $wallet->chain_id = $chainId;
        $wallet->user_id = $userId;   // set the appropriate user ID
        $wallet->favorite = false;
        $wallet->save();

        $mints = [];
        foreach ($body['tokens'] as $token) {
            $mints[] = $token['mint'];
        }

        $mintPrices = [];
        $bodyData = [
            'addresses' => $mints,
        ];
        $priceResp = $client->post(
            "/token/mainnet/prices",
            [
            'json' => $bodyData,
            ]
        );

        $priceResponse =  json_decode($priceResp->getBody(), true);
        foreach ($priceResponse as $tokenData) {
            $mint = $tokenData['tokenAddress'];
            $mintPrices[$mint] = $tokenData;
        }

        foreach ($body["tokens"] as $token) {
            $t = new Token();
            $t->name = $token["name"];
            $t->chain_id = $chainId;
            $t->current_price = $mintPrices[$token["mint"]]["usdPrice"];
            $t->address = $token["associatedTokenAddress"];
            $t->mint = $token["mint"];
            $t->symbol = $token["symbol"];
            $t->logo = $token["logo"];
            $t->save();



            //TODO: create tokenholdings here. update table aswell
            $holding = new TokenHoldings();
            $holding->user_id = $userId;
            $holding->token_id = $t->id;
            $holding->wallet_id = $wallet->id;
            $holding->amount = $token["amount"];
            $holding->value = $mintPrices[$token["mint"]]["usdPrice"] * $token["amount"];
            $walletValue +=  $mintPrices[$token["mint"]]["usdPrice"] * $token["amount"];
            $holding->save();
        }
        $wallet->value = $walletValue;
        $wallet->save();

        $snapshot = new WalletSnapshot();
        $snapshot->value = $walletValue;
        $snapshot->wallet_id = $wallet->id;
        $snapshot->save();
        return $body;

    }

    public function refreshWallet(string $walletAddress)
    {

        $wallet = Wallet::where('address', $walletAddress)->first();
        if (!$wallet) {
            throw new \Exception("Wallet with address {$walletAddress} not found.");
        }
        $wallet->refresh();
        return $wallet;
    }


    public function getTransactions(string $walletAddress)
    {
        //TODO: we get transactions and if they have tokenswap as type we add a foreign key to the transaction so we can load the swap
        todo("do this");
    }

    public function getTokenSwaps(string $walletAddress)
    {
        $wallet = DB::table("wallets")->where("address", $walletAddress)->first();
        if (!$wallet) {
            throw new \Exception("Wallet with address {$walletAddress} not found.");
        }

        $client = app(Client::class);
        $swapsResponse = $client->get(sprintf("/account/mainnet/%s/swaps?order=DESC", $walletAddress));
        $swaps = json_decode($swapsResponse->getBody(), true);

        if (!isset($swaps['result']) || empty($swaps['result'])) {
            return [];
        }

        $createdSwaps = [];

        foreach ($swaps['result'] as $swapData) {
            $tokenData = $swapData['bought']['address'] === $swapData['baseToken']
                ? $swapData['bought']
                : $swapData['sold'];

            $token = Token::where('address', $swapData['baseToken'])
                ->orWhere('mint', $tokenData['address'])
                ->first();

            if (!$token) {
                $token = Token::firstOrCreate(
                    ['address' => $swapData['baseToken']],
                    [
                        'chain_id' => 'bbdebcf5-3439-4d9c-a9e6-8e54f1924456',
                        'name' => $tokenData['name'],
                        'current_price' => $tokenData['usdPrice'],
                        'logo' => $tokenData['logo'] ?? null,
                        'symbol' => $tokenData['symbol'],
                        'mint' => $tokenData['address'],
                    ]
                );
            }

            $tokenSwap = TokenSwap::updateOrCreate(
                ['transaction_hash' => $swapData['transactionHash']],
                [
                    'chain_id' => 'bbdebcf5-3439-4d9c-a9e6-8e54f1924456',
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
                    'bought' => [
                        'address' => $swapData['bought']['address'],
                        'amount' => $swapData['bought']['amount'],
                        'usdPrice' => $swapData['bought']['usdPrice'],
                        'usdAmount' => $swapData['bought']['usdAmount'],
                        'symbol' => $swapData['bought']['symbol'],
                        'logo' => $swapData['bought']['logo'],
                        'name' => $swapData['bought']['name'],
                        'tokenType' => $swapData['bought']['tokenType'],
                    ],
                    'sold' => [
                        'address' => $swapData['sold']['address'],
                        'amount' => $swapData['sold']['amount'],
                        'usdPrice' => $swapData['sold']['usdPrice'],
                        'usdAmount' => $swapData['sold']['usdAmount'],
                        'symbol' => $swapData['sold']['symbol'],
                        'logo' => $swapData['sold']['logo'],
                        'name' => $swapData['sold']['name'],
                        'tokenType' => $swapData['sold']['tokenType'],
                    ],
                    'base_quote_price' => $swapData['baseQuotePrice'],
                    'total_value_usd' => $swapData['totalValueUsd'],
                ]
            );

            $createdSwaps[] = $tokenSwap;
        }

        return $createdSwaps;
    }

}
