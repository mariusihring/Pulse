<?php

namespace App\Models;

use App\Services\WalletService;
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

    


}
