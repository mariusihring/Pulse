<?php

namespace App\Services;

use GuzzleHttp\Client;

class WalletService
{
    public function initalWalletLoad(int $userId)
    {
        $client = app(Client::class);

        $response = $client->get("/account/mainnet/AghsmY94TE5NdCqk7FZKPW78gzQ4PpEmnKVFRao5yj9o/portfolio");
        return json_decode($response->getBody(), true);

    }
}
