<?php

namespace App\Models;

use GuzzleHttp\Client;
use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

class Wallet extends Model
{
    use HasUuids;
    use HasFactory;
    protected $fillable = [
        'address',
        'name',
        'chain_token_amount',
        'value',
        'chain_id',
        'user_id',
        'favorite'
   ];

    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    public function tokens(): HasMany
    {
        return $this->hasMany(Token::class);
    }

    public function transactions(): HasMany
    {
        return $this->hasMany(Transaction::class);
    }

    public function tokenswaps(): HasMany
    {
        return $this->hasMany(TokenSwap::class);
    }

    public function snapshots(): HasMany
    {
        return $this->hasMany(WalletSnapshot::class);
    }

    public function chain(): BelongsTo
    {
        return $this->belongsTo(Chain::class);
    }

    public function tokenholdings(): HasMany
    {
        return $this->hasMany(TokenHoldings::class);
    }

    public function refresh(): void
    {
        $client = app(Client::class);
        $response = $client->get(sprintf("/account/mainnet/%s/portfolio", $this->address));
        $body = json_decode($response->getBody(), true);

        $this->chain_token_amount = $body['nativeBalance']['solana'];
        $this->save();

        $mints = [];
        foreach ($body['tokens'] as $token) {
            $mints[] = $token['mint'];
        }

        $mintPrices = [];
        $bodyData = ['addresses' => $mints];
        $priceResp = $client->post(
            "/token/mainnet/prices",
            ['json' => $bodyData]
        );

        $priceResponse = json_decode($priceResp->getBody(), true);
        foreach ($priceResponse as $tokenData) {
            $mint = $tokenData['tokenAddress'];
            $mintPrices[$mint] = $tokenData;
        }

        $walletValue = 0;
        foreach ($body["tokens"] as $token) {
            // Update or create token with all fields
            $t = Token::updateOrCreate(
                ['mint' => $token["mint"]],
                [
                    'name' => $token["name"],
                    'chain_id' => $this->chain_id,
                    'current_price' => $mintPrices[$token["mint"]]["usdPrice"],
                    'address' => $token["associatedTokenAddress"],
                    'mint' => $token["mint"],
                    'symbol' => $token["symbol"],
                    'logo' => $token["logo"],
                ]
            );

            // Update or create token holding with all fields
            $holding = TokenHoldings::updateOrCreate(
                [
                    'user_id' => $this->user_id,
                    'token_id' => $t->id,
                    'wallet_id' => $this->id,
                ],
                [
                    'user_id' => $this->user_id,
                    'token_id' => $t->id,
                    'wallet_id' => $this->id,
                    'amount' => $token["amount"],
                    'value' => $mintPrices[$token["mint"]]["usdPrice"] * $token["amount"],
                ]
            );

            $walletValue += $mintPrices[$token["mint"]]["usdPrice"] * $token["amount"];
        }

        // Update wallet value on $this
        $this->value = $walletValue;
        $this->save();

        // Update or create wallet snapshot with all fields
        WalletSnapshot::updateOrCreate(
            ['wallet_id' => $this->id, 'created_at' => now()->startOfDay()],
            [
                'value' => $walletValue,
                'wallet_id' => $this->id,
            ]
        );
    }
}
