<?php

namespace App\Http\Controllers;

use App\Models\Chain;
use App\Models\Wallet;
use App\Services\WalletService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class WalletController extends Controller
{
    public function initialSetup(Request $request, WalletService $walletService):JsonResponse
    {
        $result = $walletService->initalWalletLoad($request->user()->id);
        return response()->json($result);
    }

    public function refreshSingleWallet(Request $request, WalletService $service) {
        $walletAddress = $request->input( "address");
        $chain = Chain::where('name', "Solana")->first();
        $data = $service->loadPortfolio(Auth::user()->id, $walletAddress, $chain->id);
        $user = Auth::user()->load([
            'tokens',
            'wallets',
            'wallets.chain',
            'wallets.snapshots',
            'wallets.tokenswaps' => function ($query) {
                $query->orderBy('block_timestamp', 'desc');
            },
            'wallets.tokenHoldings',
            'wallets.tokenHoldings.token'
        ]);

        return compact('user');
    }

public function update(Request $request, Wallet $wallet)
{
    $validated = $request->validate([
        'name' => 'sometimes|string|max:255',
        'favorite' => 'sometimes|boolean',
    ]);

    $wallet->update($validated);
    return response()->json($wallet);
}


}
