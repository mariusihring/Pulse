<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

Route::get('/user', function (Request $request) {
    return $request->user();
})->middleware('auth:sanctum');

Route::middleware(['auth', 'verified'])->get('/user/wallets', function () {
    $user = Auth::user()->load([
        'wallets',
        'wallets.chain',
        'wallets.snapshots',
        'wallets.tokenswaps' => function ($query) {
            $query->orderBy('block_timestamp', 'desc');
        },
        'wallets.tokenHoldings',
        'wallets.tokenHoldings.token',
    ]);
    return response()->json($user);
});

