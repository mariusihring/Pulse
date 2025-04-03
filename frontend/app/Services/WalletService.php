<?php

namespace App\Services;

use App\Models\Token;
use App\Models\Wallet;
use GuzzleHttp\Client;

// Maybe check here for more good apis
// https://github.com/public-apis/public-apis?tab=readme-ov-file#cryptocurrency

class WalletService
{
    public function loadPortfolio(int $userId, string $walletAddress)
    {
        $client = app(Client::class);

        $response = $client->get(sprintf("/account/mainnet/%s/portfolio", $walletAddress));
        $body =  json_decode($response->getBody(), true);

        $wallet = new Wallet;
        $wallet->address = $walletAddress;
        $wallet->name = 'My New Wallet';
        $wallet->chain_token_amount = $body['nativeBalance']['solana'];
        $wallet->value = 0;
        $wallet->chain_id = 1; // set your chain ID
        $wallet->user_id = $userId;   // set the appropriate user ID
        $wallet->favorite = false;
        $wallet->save();

        foreach ($body["tokens"] as $token) {
            $t = new Token;
            $t->name = $token["name"];
            $t->chain_id = 1;
            $t->current_price = 0;
            $t->address = $token["associatedTokenAddress"];
            $t->mint = $token["mint"];
            $t->symbol = $token["symbol"];
            $t->logo = $token["logo"];
            $t->save();



            //TODO: create tokenholdings here. update table aswell
        }

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
