<?php

namespace App\Services;

use App\Models\Token;
use App\Models\TokenHoldings;
use App\Models\Wallet;
use App\Models\WalletSnapshot;
use GuzzleHttp\Client;

// Maybe check here for more good apis
// https://github.com/public-apis/public-apis?tab=readme-ov-file#cryptocurrency

class WalletService
{
    public function loadPortfolio(string $userId, string $walletAddress)
    {
        $client = app(Client::class);
        $walletValue = 0;
        $response = $client->get(sprintf("/account/mainnet/%s/portfolio", $walletAddress));
        $body =  json_decode($response->getBody(), true);

        $wallet = new Wallet;
        $wallet->address = $walletAddress;
        $wallet->name = 'My New Wallet';
        $wallet->chain_token_amount = $body['nativeBalance']['solana'];
        $wallet->value = 0;
        $wallet->chain_id = "75f27c0a-3783-45c2-91e8-8da58222bd50";
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
        $priceResp = $client->post("/token/mainnet/prices", [
            'json' => $bodyData,
        ]);

        $priceResponse =  json_decode($priceResp->getBody(), true);
        foreach ($priceResponse as $tokenData) {
            $mint = $tokenData['tokenAddress'];
            $mintPrices[$mint] = $tokenData;
        }

        foreach ($body["tokens"] as $token) {
            $t = new Token;
            $t->name = $token["name"];
            $t->chain_id = "75f27c0a-3783-45c2-91e8-8da58222bd50";
            $t->current_price = $mintPrices[$token["mint"]]["usdPrice"];
            $t->address = $token["associatedTokenAddress"];
            $t->mint = $token["mint"];
            $t->symbol = $token["symbol"];
            $t->logo = $token["logo"];
            $t->save();



            //TODO: create tokenholdings here. update table aswell
            $holding = new TokenHoldings;
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

        $snapshot = new WalletSnapshot;
        $snapshot->value = $walletValue;
        $snapshot->wallet_id = $wallet->id;
        $snapshot->save();
        return $body;

    }

    public function refreshWallet(string $walletAddress, $lastRefresh)
    {
        todo("do this");
    }


    public function getTransfers(string $walletAddress)
    {
        todo("do this");
    }

}
