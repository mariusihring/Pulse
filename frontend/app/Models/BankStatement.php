<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class BankStatement extends Model
{
    protected $fillable = [
        'user_id',
        'my_iban',
        'receiver_iban',
        'date',
        'name_receiver',
        'usage_text',
        'amount',
        'balance_after_transaction',
    ];

    protected $casts = [
        'date' => 'date',
        'my_iban' => 'encrypted',
        'receiver_iban' => 'encrypted',
        'name_receiver' => 'encrypted',
        'usage_text' => 'encrypted',
        'amount' => 'encrypted',
        'balance_after_transaction' => 'encrypted',
    ];

    public function user()
    {
        return $this->belongsTo(User::class);
    }
}
