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
//        $data = $service->loadPortfolio("019614bf-2cac-71f9-acfb-ed2b695b57c0", "8K2MYNbuN7LvSM5132gAQXnth7oyeYLmkXcCcBm68Fbm", "7bcb8fd8-4c1f-4d09-9144-8e0ea9ca3d0e");
        $data = $service->loadPortfolio("019614bf-2cac-71f9-acfb-ed2b695b57c0", "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3", "7bcb8fd8-4c1f-4d09-9144-8e0ea9ca3d0e");
        //TODO: fetch transactions here
        dd($data);
    }
}
