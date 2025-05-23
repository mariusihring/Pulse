<?php

namespace App\Http\Controllers;

use App\Models\Chain;
use App\Models\Token;
use App\Models\TokenHoldings;
use App\Models\Wallet;
use App\Services\WalletService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Inertia\Inertia;

class CryptoController extends Controller
{
    public function dashboard()
    {
        $user = Auth::user()->load(
            [
            'wallets',
            'wallets.snapshots',
            'wallets.tokenswaps' => function ($query) {
                $query->orderBy('block_timestamp', 'desc');
            },
            'tokenHoldings',
            'tokenHoldings.token'
            ]
        );
        return Inertia::render('crypto/dashboard', compact('user'));
    }

    public function transactions()
    {
        $user = Auth::user()->load(
            [

            'wallets',
            'wallets.snapshots',


            ]
        );

        $user->tokenswaps = $user->wallets->flatMap(
            function ($wallet) {
                return $wallet->tokenswaps;
            }
        )->values();
        return Inertia::render('crypto/transactions', compact('user'));
    }

    public function tokens()
    {
        $user = Auth::user()->load(
            [
            'wallets',
            'wallets.snapshots',
            ]
        );

        $user->tokens = $user->wallets->flatMap(
            function ($wallet) {
                return $wallet->tokenswaps->pluck('token')->unique('id');
            }
        )->values();
        $user->tokenswaps = $user->wallets->flatMap(
            function ($wallet) {
                return $wallet->tokenswaps;
            }
        )->values();
        foreach ($user->wallets as $wallet) {
            foreach ($wallet->tokenswaps as $tokenswap) {
                if ($tokenswap->token) {
                    $tokenswap->token->pnl = $tokenswap->token->calculatePnL();
                }
            }
        }

        return Inertia::render('crypto/tokens', compact('user'));
    }

    public function wallets()
    {
        $user = Auth::user()->load(
            [
            'wallets',
            'wallets.chain',
            'wallets.snapshots',
            'wallets.tokenswaps' => function ($query) {
                $query->orderBy('block_timestamp', 'desc');
            },
            'wallets.tokenHoldings',
            'wallets.tokenHoldings.token'
            ]
        );

        return Inertia::render('crypto/wallets', compact('user'));
    }


    public function userWallets(Request $request)
    {
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

    public function updatePortfolio(Request $request, WalletService $service)
    {
        $wallets = Wallet::where('user_id', Auth::user()->id)->get();
        $chain = Chain::where('name', "Solana")->first();
        foreach ($wallets as $wallet) {
            $service->loadPortfolio(Auth::user()->id, $wallet->address, $chain->id);
        }

        $user = Auth::user()->load(
            [
            'tokens',
            'wallets',
            'wallets.chain',
            'wallets.snapshots',
            'wallets.tokenswaps' => function ($query) {
                $query->orderBy('block_timestamp', 'desc');
            },
            'wallets.tokenHoldings',
            'wallets.tokenHoldings.token'
            ]
        );

        return compact('user');
    }

}
