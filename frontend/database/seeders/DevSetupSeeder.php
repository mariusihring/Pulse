<?php

namespace Database\Seeders;

use App\Models\Chain;
use App\Models\User;
use App\Services\WalletService;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class DevSetupSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $user = new User;
        $user->name = "Riri";
        $user->email = "test@test.de";
        $user->password = Hash::make("password");
        $user->save();


        $chain = new Chain;
        $chain->name = 'Solana';
        $chain->save();


        //user wallet chain
        $service = new WalletService;
        $data1 = $service->loadPortfolio($user->id, "8K2MYNbuN7LvSM5132gAQXnth7oyeYLmkXcCcBm68Fbm", $chain->id);
        $data1 = $service->loadPortfolio($user->id, "AghsmY94TE5NdCqk7FZKPW78gzQ4PpEmnKVFRao5yj9o", $chain->id);
        $data2 = $service->loadPortfolio($user->id, "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3", $chain->id);
        $data3 = $service->getTokenSwaps("AghsmY94TE5NdCqk7FZKPW78gzQ4PpEmnKVFRao5yj9o", $chain->id);

    }
}
