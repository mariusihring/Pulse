<?php

namespace App\Http\Controllers;

use App\Services\WalletService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class WalletController extends Controller
{
    public function initialSetup(Request $request, WalletService $walletService):JsonResponse
    {
        $result = $walletService->initalWalletLoad($request->user()->id);
        return response()->json($result);
    }
}
