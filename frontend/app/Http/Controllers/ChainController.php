<?php

namespace App\Http\Controllers;

use App\Models\Chain;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;

class ChainController extends Controller
{
    /**
     * Retrieve all available blockchain chains.
     *
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        try {
            $chains = Chain::select('id', 'name')->get();

            return response()->json([
                'success' => true,
                'chains' => $chains,
            ], 200);
        } catch (\Exception $e) {
            return response()->json([
                'success' => false,
                'message' => 'Failed to fetch chains: ' . $e->getMessage(),
            ], 500);
        }
    }
}
