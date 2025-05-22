<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

class Chain extends Model
{

    use HasFactory, HasUuids;
    protected $fillable = [
        "name"
    ];

    public function tokens(): HasMany
    {
        return $this->hasMany(Token::class);
    }


    public function tokenSwaps(): HasMany
    {
        return $this->hasMany(TokenSwap::class);
    }

    public function refreshChainTokenPrice(): void
    {
        todo("implement this with moralis api");
    }
}
