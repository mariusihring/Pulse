<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;

class WalletSnapshot extends Model
{
    use HasUuids;

    protected $fillable = [
        "wallet_id",
        "value",

    ];
}
