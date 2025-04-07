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
            $table->uuid('id')->primary()->default(DB::raw('gen_random_uuid()'));
            $table->string('address')->unique();
            $table->string('name')->nullable();
            $table->double("chain_token_amount");
            $table->double('value');
            $table->foreignUuid('chain_id');
            $table->foreignUuid('user_id');
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
