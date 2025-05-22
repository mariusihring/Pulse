<?php

namespace App\Http\Controllers;

use App\Services\BankStatementService;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Redirect;
use App\Models\BankStatement;


class BankingController extends Controller
{
    protected $bankStatementService;

    public function __construct(BankStatementService $bankStatementService)
    {
        $this->bankStatementService = $bankStatementService;
    }

    public function index()
    {
        $user = Auth::user();
        $statements = BankStatement::where('user_id', $user->id)
            ->orderBy('date', 'desc')
            ->get()
            ->map(function ($statement) {
                return [
                    'id' => $statement->id,
                    'my_iban' => $statement->my_iban,
                    'receiver_iban' => $statement->receiver_iban,
                    'date' => $statement->date->format('Y-m-d'),
                    'name_receiver' => $statement->name_receiver,
                    'usage_text' => $statement->usage_text,
                    'amount' => $statement->amount,
                    'balance_after_transaction' => $statement->balance_after_transaction,
                ];
            });

        return Inertia::render('banking/dashboard', compact('statements'));
    }

    public function uploadCsv(Request $request)
    {
        $request->validate([
            'csv_file' => 'required|file|mimes:csv,txt|max:2048', // Max 2MB
        ]);

        try {
            $user = Auth::user();
            $results = $this->bankStatementService->processCsvUpload($request->file('csv_file'), $user->id);

            return Redirect::route('bank.statements')->with([
                'success' => "Successfully processed {$results['success']} rows.",
                'errors' => $results['errors'],
            ]);
        } catch (\Exception $e) {
            return Redirect::route('bank.statements')->withErrors(['csv_file' => $e->getMessage()]);
        }
    }
}
