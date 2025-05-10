<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use Inertia\Inertia;

Route::get(
    '/',
    function () {
        return Inertia::render('welcome');
    }
)->name('home');

Route::middleware(['auth', 'verified'])->group(
    function () {
        Route::get(
            'dashboard',
            function () {
                return Inertia::render('dashboard');
            }
        )->name('dashboard');

        Route::get('/bank/statements', [\App\Http\Controllers\BankingController::class, 'index'])->middleware('auth')->name('bank.statements');
        Route::post('/bank/statements/upload', [\App\Http\Controllers\BankingController::class, 'uploadCsv'])->middleware('auth')->name('bank.statements.upload');
    }
);

require __DIR__ . '/settings.php';
require __DIR__ . '/auth.php';
require __DIR__ . '/crypto.php';
//require __DIR__ . '/banking.php';
