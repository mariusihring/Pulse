<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Token extends Model
{
    //

    public function chain(): BelongsTo
    {
        return $this->belongsTo(Chain::class);
    }
}
