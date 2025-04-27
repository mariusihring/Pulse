<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use Illuminate\Support\Facades\Auth;
use Inertia\Inertia;

Route::middleware(['auth', 'verified'])->group(function () {
    Route::prefix('crypto')->group(function () {
        Route::get('/dashboard', function () {
            $user = Auth::user()->load([
                'wallets',
                'wallets.snapshots',
                'wallets.tokenswaps' => function ($query) {
                    $query->orderBy('block_timestamp', 'desc');
                },
                'tokenHoldings',
                'tokenHoldings.token'
            ]);
            return Inertia::render('crypto/dashboard', compact('user'));
        })->name('cryptodashboard');

        Route::get('/transactions', function () {
            $user = Auth::user()->load([
                'wallets',
                'wallets.snapshots',
                'wallets.tokenswaps' => function ($query) {
                    $query->orderBy('block_timestamp', 'desc');
                },
                'tokenHoldings',
                'tokenHoldings.token'
            ]);
            return Inertia::render('crypto/transactions', compact('user'));
        })->name('crypto.transactions');

        Route::post('/test', function (Request $request, \App\Services\WalletService $service) {
            $user = Auth::user();
            $address = $request->input('address');
            $data = $service->loadPortfolio($user->id, $address, 'bbdebcf5-3439-4d9c-a9e6-8e54f1924456');
            return response()->json($data);
        })->name('crypto.test');

        Route::post('/refresh', function (Request $request, \App\Services\WalletService $service) {
            $user = Auth::user();
            $address = $request->input('address');
            $data = $service->refreshWallet($address);
            return response()->json($data);
        })->name('crypto.refresh');
    });
}); 