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
    public function run(WalletService $service): void
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
        $data1 = $service->loadPortfolio($user->id, config('portfolio.ids.portfolio_1'), $chain->id);
        $data1 = $service->loadPortfolio($user->id, config('portfolio.ids.portfolio_2'), $chain->id);
        $data2 = $service->loadPortfolio($user->id, config('portfolio.ids.portfolio_3'), $chain->id);
        $data3 = $service->getTokenSwaps(config('portfolio.ids.portfolio_2'), $chain->id);
        $data3 = $service->getTokenSwaps(config('portfolio.ids.portfolio_1'), $chain->id);
        $data3 = $service->getTokenSwaps(config('portfolio.ids.portfolio_3'), $chain->id);

    }
}
