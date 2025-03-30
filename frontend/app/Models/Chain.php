<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

class Chain extends Model
{
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
