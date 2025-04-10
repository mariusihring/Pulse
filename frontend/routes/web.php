<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use Inertia\Inertia;

Route::get('/', function () {
    return Inertia::render('welcome');
})->name('home');

Route::middleware(['auth', 'verified'])->group(function () {
    Route::get('dashboard', function () {

        $user = Auth::user()->load([
                'wallets',
                'wallets.snapshots',
                'wallets.tokenswaps',
                'tokenHoldings',
                'tokenHoldings.token'
            ]
        );
        return Inertia::render('dashboard', compact('user'));
    })->name('dashboard');

    Route::post("test" ,function (Request $request, \App\Services\WalletService $service) {
        $user = Auth::user();
        $address = $request->input('address');
        $data = $service->loadPortfolio($user->id, $address, 'bbdebcf5-3439-4d9c-a9e6-8e54f1924456');

        return response()->json($data);

    })->name('test');
});

require __DIR__ . '/settings.php';
require __DIR__ . '/auth.php';
