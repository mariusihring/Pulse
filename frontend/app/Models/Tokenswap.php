<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class TokenSwap extends Model
{
    use HasFactory, HasUuids;

    /**
     * The table associated with the model.
     *
     * @var string
     */
    protected $table = 'tokenswaps';

    /**
     * The primary key for the model.
     *
     * @var string
     */
    protected $primaryKey = 'id';

    /**
     * Indicates if the IDs are auto-incrementing.
     *
     * @var bool
     */
    public $incrementing = false;

    /**
     * The "type" of the primary key ID.
     *
     * @var string
     */
    protected $keyType = 'string';

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'chain_id',
        'token_id',
        'wallet_id',
        'transaction_hash',
        'transaction_type',
        'transaction_index',
        'sub_category',
        'block_timestamp',
        'block_number',
        'wallet_address',
        'pair_address',
        'pair_label',
        'exchange_address',
        'exchange_name',
        'exchange_logo',
        'base_token',
        'quote_token',
        'bought',
        'sold',
        'base_quote_price',
        'total_value_usd',
    ];

    /**
     * The attributes that should be cast.
     *
     * @var array
     */
    protected $casts = [
        'id' => 'string',
        'chain_id' => 'string',
        'token_id' => 'string',
        'wallet_id' => 'string',
        'block_timestamp' => 'datetime',
        'bought' => 'array', // JSON cast to array
        'sold' => 'array',   // JSON cast to array
        'base_quote_price' => 'decimal:18',
        'total_value_usd' => 'decimal:8',
    ];

    /**
     * The attributes that should be hidden for serialization.
     *
     * @var array
     */
    protected $hidden = [
        'created_at',
        'updated_at',
    ];

    /**
     * Get the chain that this token swap belongs to.
     */
    public function chain(): BelongsTo
    {
        return $this->belongsTo(Chain::class, 'chain_id');
    }

    /**
     * Get the token involved in this swap.
     */
    public function token(): BelongsTo
    {
        return $this->belongsTo(Token::class, 'token_id');
    }

    /**
     * Get the wallet that performed this swap.
     */
    public function wallet(): BelongsTo
    {
        return $this->belongsTo(Wallet::class, 'wallet_id');
    }
}
