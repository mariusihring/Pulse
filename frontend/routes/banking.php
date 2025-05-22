<?php

use App\Http\Controllers\CryptoController;
use Illuminate\Support\Facades\Route;


Route::middleware(['auth', 'verified'])->group(function () {
    Route::prefix('banking')->group(function () {
        Route::get('/csv-upload', [\App\Http\Controllers\BankingController::class, 'csvUpload'])->name('bankind.csv.upload');
        Route::post( "/upload", [\App\Http\Controllers\BankingController::class,  'processCsvUpload'])->name( "banking.csv.upload.processing");
    });
});
