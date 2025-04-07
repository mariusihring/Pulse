<?php

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
                'tokenHoldings',
                'tokenHoldings.token'
            ]
        );
        return Inertia::render('dashboard', compact('user'));
    })->name('dashboard');
});

require __DIR__ . '/settings.php';
require __DIR__ . '/auth.php';
