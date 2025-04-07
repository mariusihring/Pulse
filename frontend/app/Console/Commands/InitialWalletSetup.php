<?php

namespace App\Console\Commands;

use App\Services\WalletService;
use GuzzleHttp\Client;
use Illuminate\Console\Command;

class InitialWalletSetup extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'app:initial-wallet-setup';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Command description';

    /**
     * Execute the console command.
     */
    public function handle(WalletService $service)
    {
        $data = $service->loadPortfolio("0196100a-9598-72c7-923e-a42c7ac61ff7", "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3");
        //TODO: fetch transactions here
        dd($data);
    }
}
