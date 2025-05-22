<?php

namespace App\Services;

use App\Models\BankStatement;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class BankStatementService
{
    public function processCsvUpload($file, string $userId): array
    {
        try {
            $path = $file->getRealPath();
            $handle = fopen($path, 'r');
            if ($handle === false) {
                throw new \Exception('Failed to open CSV file.');
            }

            // Read header with semicolon delimiter
            $header = fgetcsv($handle, 1000, ';');
            $requiredColumns = [
                'Bezeichnung Auftragskonto',
                'IBAN Auftragskonto',
                'BIC Auftragskonto',
                'Bankname Auftragskonto',
                'Buchungstag',
                'Valutadatum',
                'Name Zahlungsbeteiligter',
                'IBAN Zahlungsbeteiligter',
                'BIC (SWIFT-Code) Zahlungsbeteiligter',
                'Buchungstext',
                'Verwendungszweck',
                'Betrag',
                'Waehrung',
                'Saldo nach Buchung',
                'Bemerkung',
                'Kategorie',
                'Steuerrelevant',
                'Glaeubiger ID',
                'Mandatsreferenz',
            ];
            if (!$header || array_diff($requiredColumns, $header)) {
                throw new \Exception('Invalid CSV format. Required columns: ' . implode(', ', $requiredColumns));
            }

            $results = ['success' => 0, 'errors' => []];
            DB::beginTransaction();

            while (($row = fgetcsv($handle, 1000, ';')) !== false) {
                $data = array_combine($header, $row);
                try {
                    // Validate data
                    if (empty($data['Buchungstag']) || !strtotime($data['Buchungstag']) ||
                        empty($data['IBAN Auftragskonto']) || empty($data['IBAN Zahlungsbeteiligter']) ||
                        empty($data['Name Zahlungsbeteiligter']) || empty($data['Verwendungszweck']) ||
                        !is_numeric(str_replace(',', '.', $data['Betrag'])) ||
                        !is_numeric(str_replace(',', '.', $data['Saldo nach Buchung']))) {
                        throw new \Exception('Missing or invalid data in row.');
                    }

                    // Create bank statement
                    BankStatement::create([
                        'user_id' => $userId,
                        'my_iban' => $data['IBAN Auftragskonto'],
                        'receiver_iban' => $data['IBAN Zahlungsbeteiligter'],
                        'date' => date('Y-m-d', strtotime($data['Buchungstag'])),
                        'name_receiver' => $data['Name Zahlungsbeteiligter'],
                        'usage_text' => $data['Verwendungszweck'],
                        'amount' => str_replace(',', '.', $data['Betrag']),
                        'balance_after_transaction' => str_replace(',', '.', $data['Saldo nach Buchung']),
                    ]);

                    $results['success']++;
                } catch (\Exception $e) {
                    $results['errors'][] = "Error in row: {$e->getMessage()}";
                    Log::error("CSV processing error", ['row' => $data, 'error' => $e->getMessage()]);
                }
            }

            fclose($handle);
            DB::commit();
            return $results;
        } catch (\Exception $e) {
            DB::rollBack();
            Log::error("CSV upload failed: {$e->getMessage()}");
            throw $e;
        }
    }
}
