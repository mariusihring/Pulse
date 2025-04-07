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
        Schema::create('tokenswaps', function (Blueprint $table) {
            $table->uuid('id')->primary()->default(DB::raw('gen_random_uuid()'));
            $table->foreignUuid('chain_id')->constrained();
            $table->foreignUuid('token_id')->constrained();
            $table->foreignUuid('wallet_id')->constrained();
            $table->string('transaction_hash')->unique();
            $table->string('transaction_type');
            $table->integer('transaction_index');
            $table->string('sub_category')->nullable();
            $table->timestamp('block_timestamp');
            $table->unsignedBigInteger('block_number');
            $table->string('wallet_address');
            $table->string('pair_address');
            $table->string('pair_label');
            $table->string('exchange_address');
            $table->string('exchange_name');
            $table->string('exchange_logo')->nullable();
            $table->string('base_token');
            $table->string('quote_token');
            $table->json('bought');
            $table->json('sold');
            $table->decimal('base_quote_price', 30, 18);
            $table->decimal('total_value_usd', 20, 8);
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('tokenswaps');
    }
};
