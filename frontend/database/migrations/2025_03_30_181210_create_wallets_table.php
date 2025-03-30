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
        Schema::create('wallets', function (Blueprint $table)
        {
            $table->id();
            $table->string('address')->unique();
            $table->string('name')->nullable();
            $table->double("chain_token_amount");
            $table->double('value');
            $table->foreignId('chain_id');
            $table->foreignId('user_id');
            $table->boolean('favorite')->default(false);
            $table->timestamps();
        }
        );
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('wallets');
    }
};
