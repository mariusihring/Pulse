<?php

use App\Http\Controllers\CryptoController;
use Illuminate\Support\Facades\Route;


Route::middleware(['auth', 'verified'])->group(function () {
    Route::prefix('crypto')->group(function () {
        Route::get('/dashboard', [CryptoController::class, 'dashboard'])->name('cryptodashboard');
        Route::get('/transactions', [CryptoController::class, 'transactions'])->name('crypto.transactions');
        Route::get( "/tokens", [CryptoController::class,  'tokens'])->name("crypto.tokens");
        Route::get('/wallets', [CryptoController::class, 'wallets'])->name('crypto.wallets');
        Route::post('/test', [CryptoController::class, 'test'])->name('crypto.test');
        Route::get('/user/wallets/all', [CryptoController::class, 'userWallets'])->name('crypto.user.wallet.reload');
        Route::post("/user/wallets/refresh", [CryptoController::class, 'updatePortfolio'])->name('crypto.user.wallets.updatePortfolio');
        Route::post("/user/wallet/refresh", [\App\Http\Controllers\WalletController::class, 'refreshSingleWallet'])->name('crypto.user.wallets.refresh.single.wallet');
        Route::patch("/user/wallets/{wallet}", [\App\Http\Controllers\WalletController::class,  'update'])->name( 'crypto.user.wallet.update');
        Route::get( "/chains", [\App\Http\Controllers\ChainController::class,  'index'])->name("crypto.chains.all");
    });
});
