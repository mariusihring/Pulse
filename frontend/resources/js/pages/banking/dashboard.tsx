import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { useState } from 'react';
import { useForm } from '@inertiajs/react'
import {toast} from "sonner"

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Banking Dashboard',
        href: '/banking/dashboard',
    },
];
interface BankStatement {
    id: string;
    my_iban: string;
    receiver_iban: string;
    date: string;
    name_receiver: string;
    usage_text: string;
    amount: string; // Decrypted as string due to encryption
    balance_after_transaction: string;
  }

export default function Dashboard({ statements, errors, success, errors_array}: {statements: BankStatement[]}) {

    const [searchQuery, setSearchQuery] = useState('');

    // Form for CSV upload
    const { data, setData, post, progress, processing } = useForm({
      csv_file: null as File | null,
    });

    // Handle file change
    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files && e.target.files[0]) {
        setData('csv_file', e.target.files[0]);
      }
    };

    // Handle form submission
    const handleSubmit = (e: React.FormEvent) => {
      e.preventDefault();
      post(route('bank.statements.upload'), {
        onSuccess: () => {
          toast.success(success || 'CSV uploaded successfully!');
          setData('csv_file', null);
        },
        onError: () => {
          toast.error(errors?.csv_file || 'Failed to upload CSV.');
        },
      });
    };

    // Filter statements
    const filteredStatements = statements.filter(
      (statement) =>
        statement.name_receiver.toLowerCase().includes(searchQuery.toLowerCase()) ||
        statement.usage_text.toLowerCase().includes(searchQuery.toLowerCase())
    );

    // Handle refresh
    const handleRefresh = () => {
      Inertia.reload({
        onSuccess: () => toast.success('Statements refreshed successfully!'),
        onError: () => toast.error('Failed to refresh statements.'),
      });
    };
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Banking Upload" />
             <div className="container mx-auto p-6 max-w-7xl">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">Bank Statements</h1>

      {/* CSV Upload Form */}
      <div className="mb-6 bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-semibold mb-4">Upload Bank Statement CSV</h2>
        <form onSubmit={handleSubmit} encType="multipart/form-data">
          <div className="flex flex-col sm:flex-row gap-4">
            <input
              type="file"
              accept=".csv"
              onChange={handleFileChange}
              className="border border-gray-300 rounded-lg p-2 w-full sm:w-1/2"
            />
            <button
              type="submit"
              disabled={processing || !data.csv_file}
              className="bg-blue-500 text-white rounded-lg px-4 py-2 hover:bg-blue-600 transition disabled:bg-gray-400"
            >
              {processing ? 'Uploading...' : 'Upload CSV'}
            </button>
          </div>
          {progress && (
            <div className="mt-2">
              <progress value={progress.percentage} max="100" className="w-full">
                {progress.percentage}%
              </progress>
            </div>
          )}
          {errors?.csv_file && <p className="text-red-500 mt-2">{errors.csv_file}</p>}
          {success && <p className="text-green-500 mt-2">{success}</p>}
          {errors_array && errors_array.length > 0 && (
            <div className="mt-2">
              <p className="text-red-500">Processing errors:</p>
              <ul className="list-disc pl-5 text-red-500">
                {errors_array.map((error, index) => (
                  <li key={index}>{error}</li>
                ))}
              </ul>
            </div>
          )}
        </form>
        <p className="text-gray-500 text-sm mt-2">
          CSV must include columns: Bezeichnung Auftragskonto, IBAN Auftragskonto, BIC Auftragskonto, Bankname Auftragskonto, Buchungstag, Valutadatum, Name Zahlungsbeteiligter, IBAN Zahlungsbeteiligter, BIC (SWIFT-Code) Zahlungsbeteiligter, Buchungstext, Verwendungszweck, Betrag, Waehrung, Saldo nach Buchung, Bemerkung, Kategorie, Steuerrelevant, Glaeubiger ID, Mandatsreferenz.
        </p>
      </div>

      {/* Search and Refresh */}
      <div className="flex flex-col sm:flex-row justify-between mb-6 gap-4">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search by receiver name or usage text..."
          className="border border-gray-300 rounded-lg p-2 w-full sm:w-1/3 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button
          onClick={handleRefresh}
          className="bg-blue-500 text-white rounded-lg px-4 py-2 hover:bg-blue-600 transition"
        >
          Refresh
        </button>
      </div>

      {/* Statements Table */}
      {filteredStatements.length === 0 ? (
        <p className="text-gray-500 text-center">No statements found.</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full bg-white border border-gray-200 rounded-lg shadow-md">
            <thead>
              <tr className="bg-gray-100">
                <th className="border-b p-4 text-left text-gray-600">Date</th>
                <th className="border-b p-4 text-left text-gray-600">My IBAN</th>
                <th className="border-b p-4 text-left text-gray-600">Receiver IBAN</th>
                <th className="border-b p-4 text-left text-gray-600">Receiver Name</th>
                <th className="border-b p-4 text-left text-gray-600">Usage Text</th>
                <th className="border-b p-4 text-left text-gray-600">Amount</th>
                <th className="border-b p-4 text-left text-gray-600">Balance After</th>
              </tr>
            </thead>
            <tbody>
              {filteredStatements.map((statement) => (
                <tr key={statement.id} className="hover:bg-gray-50">
                  <td className="border-b p-4">{statement.date}</td>
                  <td className="border-b p-4">
                    {statement.my_iban.slice(-4).padStart(statement.my_iban.length, '*')}
                  </td>
                  <td className="border-b p-4">
                    {statement.receiver_iban.slice(-4).padStart(statement.receiver_iban.length, '*')}
                  </td>
                  <td className="border-b p-4">{statement.name_receiver}</td>
                  <td className="border-b p-4">{statement.usage_text}</td>
                  <td className="border-b p-4">€{parseFloat(statement.amount).toFixed(2)}</td>
                  <td className="border-b p-4">€{parseFloat(statement.balance_after_transaction).toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}</div>
        </AppLayout>
    );
}
