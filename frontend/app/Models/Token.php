<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Concerns\HasVersion4Uuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

use function GuzzleHttp\json_encode;

class Token extends Model
{
    use HasFactory;
    use HasUuids;

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'chain_id',
        'name',
        'current_price',
        'logo',
        'symbol',
        'address',
        'mint',
    ];

    /**
     * Get the chain that owns the token.
     */
    public function chain()
    {
        return $this->belongsTo(Chain::class);
    }

    public function tokenswaps()
    {
        return $this->hasMany(TokenSwap::class);
    }

    public function calculatePnL()
    {
        $tokenswaps = $this->tokenswaps;
        $totalBought = 0;
        $totalSold = 0;
        $totalBoughtAmount = 0;
        $totalSoldAmount = 0;

        \Log::info('Calculating PnL for token: ' . $this->symbol);
        \Log::info('Current price: ' . $this->current_price);

        foreach ($tokenswaps as $tokenswap) {
            \Log::info('Processing swap:', [
                'type' => $tokenswap->transaction_type,
                'bought' => $tokenswap->bought,
                'sold' => $tokenswap->sold
            ]);

            if ($tokenswap->transaction_type === 'buy') {
                $totalBought += $tokenswap->bought['usdAmount'];
                $totalBoughtAmount += $tokenswap->bought['amount'];
            } else {
                $totalSold += $tokenswap->sold['usdAmount'];
                $totalSoldAmount += $tokenswap->sold['amount'];
            }
        }

        \Log::info('Totals:', [
            'totalBought' => $totalBought,
            'totalSold' => $totalSold,
            'totalBoughtAmount' => $totalBoughtAmount,
            'totalSoldAmount' => $totalSoldAmount
        ]);

        $realizedPnL = 0;
        $unrealizedPnL = 0;

        if ($totalBoughtAmount > 0) {
            $averageBuyPrice = $totalBought / $totalBoughtAmount;

            if ($totalSoldAmount > 0) {
                $realizedPnL = $totalSold - ($totalSoldAmount * $averageBuyPrice);
            }

            $remainingAmount = $totalBoughtAmount - $totalSoldAmount;
            if ($remainingAmount > 0) {
                $unrealizedPnL = ($this->current_price - $averageBuyPrice) * $remainingAmount;
            }
        }

        $totalPnL = $realizedPnL + $unrealizedPnL;

        return
            [
            'totalBoughtUsd' => $totalBought,
            'totalSoldUsd' => $totalSold,
            'totalTokenAmountBought' => $totalBoughtAmount,
            'totalTokenAmountSold' => $totalSoldAmount,
            'realizedPnL' => $realizedPnL,
            'unrealizedPnL' => $unrealizedPnL,
            'totalPnL' => $totalPnL,
            'averageBuyPrice' => $totalBoughtAmount > 0 ? $totalBought / $totalBoughtAmount : 0,
            'remainingAmount' => $totalBoughtAmount - $totalSoldAmount
            ]
        ;

    }

}
