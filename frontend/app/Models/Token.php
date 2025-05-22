<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Concerns\HasVersion4Uuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Token extends Model
{
    use HasFactory, HasUuids;

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

            if ($tokenswap->transaction_type === 'Buy') {
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

        // Initialize PnL values
        $realizedPnL = 0;
        $unrealizedPnL = 0;
        
        // Only calculate if we have bought tokens
        if ($totalBoughtAmount > 0) {
            $averageBuyPrice = $totalBought / $totalBoughtAmount;
            
            // Calculate realized PnL from completed trades
            if ($totalSoldAmount > 0) {
                $realizedPnL = $totalSold - ($totalSoldAmount * $averageBuyPrice);
            }
            
            // Calculate unrealized PnL from remaining tokens
            $remainingAmount = $totalBoughtAmount - $totalSoldAmount;
            if ($remainingAmount > 0) {
                $unrealizedPnL = ($this->current_price - $averageBuyPrice) * $remainingAmount;
            }
        }
        
        // Total PnL is sum of realized and unrealized
        $totalPnL = $realizedPnL + $unrealizedPnL;
        
        \Log::info('Final calculation:', [
            'realizedPnL' => $realizedPnL,
            'unrealizedPnL' => $unrealizedPnL,
            'totalPnL' => $totalPnL,
            'averageBuyPrice' => $totalBoughtAmount > 0 ? $totalBought / $totalBoughtAmount : 0,
            'remainingAmount' => $totalBoughtAmount - $totalSoldAmount
        ]);

        return $totalPnL;
    }
    
}
