<?php

use Illuminate\Foundation\Inspiring;
use Illuminate\Support\Facades\Artisan;
use Illuminate\Support\Facades\Schedule;
use App\Jobs\RefreshAllWallets;



Schedule::call(function () {
    \Log::info("Dispatching wallet refresh job");
    dispatch(new RefreshAllWallets);
})->daily();
