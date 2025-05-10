<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('bank_statements', function (Blueprint $table) {
            $table->uuid( 'id')->primary();
            $table->foreignUuid('user_id')->constrained()->onDelete('cascade');
            $table->text('my_iban'); // Encrypted
            $table->text('receiver_iban'); // Encrypted
            $table->date('date');
            $table->text('name_receiver'); // Encrypted
            $table->text('usage_text'); // Encrypted
            $table->text('amount'); // Encrypted
            $table->text('balance_after_transaction'); // Encrypted
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('bank_statements');
    }
};
