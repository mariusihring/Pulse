<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Tokenswap extends Model
{
    use HasUuids, HasFactory;

    function wallet(): BelongsTo
    {
        return $this->belongsTo(Wallet::class);
    }
}
