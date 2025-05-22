<?php

namespace App\Jobs;

use App\Models\Wallet;
use App\Services\WalletService;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Log;

class RefreshAllWallets implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * The number of times the job may be attempted.
     */
    public $tries = 3;

    /**
     * The number of seconds to wait before retrying the job.
     */
    public $backoff = 60;

    /**
     * Get the tags that should be assigned to the job.
     */
    public function tags(): array
    {
        return ['wallet-refresh', 'crypto'];
    }

    /**
     * Execute the job.
     */
    public function handle(WalletService $walletService): void
    {
        Log::info('Starting wallet refresh job');
        
        $wallets = Wallet::all();
        $totalWallets = $wallets->count();
        $successfulRefreshes = 0;
        $failedRefreshes = 0;
        
        Log::info("Found {$totalWallets} wallets to refresh");
        
        foreach ($wallets as $wallet) {
            try {
                $walletService->refreshWallet($wallet->address);
                $successfulRefreshes++;
                Log::info("Successfully refreshed wallet: {$wallet->address}");
            } catch (\Exception $e) {
                $failedRefreshes++;
                Log::error("Failed to refresh wallet {$wallet->address}: " . $e->getMessage());
                throw $e; // Re-throw to allow Horizon to handle retries
            }
        }
        
        Log::info("Wallet refresh job completed. Success: {$successfulRefreshes}, Failed: {$failedRefreshes}, Total: {$totalWallets}");
    }

    /**
     * Handle a job failure.
     */
    public function failed(\Throwable $exception): void
    {
        Log::error('Wallet refresh job failed: ' . $exception->getMessage());
    }
} 