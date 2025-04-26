<?php

namespace App\Console\Commands;

use App\Services\WalletService;
use GuzzleHttp\Client;
use Illuminate\Console\Command;
use function GuzzleHttp\json_encode;

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
        $data1 = $service->loadPortfolio("01961ec3-cfc9-7395-947c-5517b9404507", "8K2MYNbuN7LvSM5132gAQXnth7oyeYLmkXcCcBm68Fbm", "bbdebcf5-3439-4d9c-a9e6-8e54f1924456");
        $data2 = $service->loadPortfolio("01961ec3-cfc9-7395-947c-5517b9404507", "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3", "bbdebcf5-3439-4d9c-a9e6-8e54f1924456");
        $data3 = $service->getTokenSwaps("AghsmY94TE5NdCqk7FZKPW78gzQ4PpEmnKVFRao5yj9o");
        //TODO: fetch transactions here
        dd(json_encode([$data1, $data2, $data3]));
    }
}
