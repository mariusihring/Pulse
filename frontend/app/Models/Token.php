<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Token extends Model
{
    use HasFactory;

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
}
