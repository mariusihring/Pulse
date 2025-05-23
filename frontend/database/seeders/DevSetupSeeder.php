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
        $user = new User();
        $user->name = "Riri";
        $user->email = "test@test.de";
        $user->password = Hash::make("password");
        $user->save();


        $chain = new Chain();
        $chain->name = 'Solana';
        $chain->save();

        $portfolio1 = env('PORTFOLIO_ID_1');
        $portfolio2 = env('PORTFOLIO_ID_2');
        $portfolio3 = env('PORTFOLIO_ID_3');

        //user wallet chain
        $data1 = $service->loadPortfolio($user->id, $portfolio1, $chain->id);
        $data1 = $service->loadPortfolio($user->id, $portfolio2, $chain->id);
        $data1 = $service->loadPortfolio($user->id, $portfolio3, $chain->id);
        $data3 = $service->getTokenSwaps($portfolio1, $chain->id);
        $data3 = $service->getTokenSwaps($portfolio2, $chain->id);
        $data3 = $service->getTokenSwaps($portfolio3, $chain->id);

    }
}
