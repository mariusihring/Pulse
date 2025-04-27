import { Breadcrumbs } from '@/components/breadcrumbs';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { type BreadcrumbItem as BreadcrumbItemType } from '@/types';
import { Button } from '@/components/ui/button';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useState } from 'react';
import axios from 'axios';
import { RefreshCcw } from 'lucide-react';

export function AppSidebarHeader({ breadcrumbs = [] }: { breadcrumbs?: BreadcrumbItemType[] }) {
    const [isOpen, setIsOpen] = useState(false); // Control Popover open state
    const [address, setAddress] = useState('');  // Store input address
    const [responseData, setResponseData] = useState(''); // Store API response
    const [loading, setLoading] = useState(false); // Handle loading state
    const [error, setError] = useState(''); // Handle errors

    // Handle form submission
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setResponseData('');

        try {
            const response = await axios.post(
                '/test', // Adjust if your API URL differs (e.g., '/api/test')
                { address: address },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]')?.getAttribute('content'), // For Laravel CSRF
                    },
                }
            );
            setResponseData(JSON.stringify(response.data, null, 2)); // Pretty-print JSON response
        } catch (err) {
            setError('Failed to load portfolio data. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };


    // Handle form submission
    const handleRefresh = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setResponseData('');

        try {
            const response = await axios.post(
                '/refresh', // Adjust if your API URL differs (e.g., '/api/test')
                { address: "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3" },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]')?.getAttribute('content'), // For Laravel CSRF
                    },
                }
            );
            console.log(response.data)
        } catch (err) {
            //setError('Failed to load portfolio data. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };


    return (
        <header className="border-sidebar-border/50 flex h-16 shrink-0 items-center gap-2 border-b px-6 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 md:px-4">
            <div className="flex items-center gap-2">
                <SidebarTrigger className="-ml-1" />
                <Breadcrumbs breadcrumbs={breadcrumbs} />
            </div>
            {/*<div className="flex w-full justify-end space-x-4">*/}
            {/*    <Button onClick={handleRefresh}><RefreshCcw /></Button>*/}
            {/*    <Popover open={isOpen} onOpenChange={setIsOpen}>*/}
            {/*        <PopoverTrigger asChild>*/}
            {/*            <Button onClick={() => setIsOpen(true)}>+</Button>*/}
            {/*        </PopoverTrigger>*/}
            {/*        <PopoverContent className="w-80">*/}
            {/*            <form onSubmit={handleSubmit} className="space-y-4">*/}
            {/*                <div>*/}
            {/*                    <label htmlFor="address" className="block text-sm font-medium text-gray-700">*/}
            {/*                        Wallet Address*/}
            {/*                    </label>*/}
            {/*                    <Input*/}
            {/*                        id="address"*/}
            {/*                        value={address}*/}
            {/*                        onChange={(e) => setAddress(e.target.value)}*/}
            {/*                        placeholder="Enter wallet address"*/}
            {/*                        className="mt-1"*/}
            {/*                        disabled={loading}*/}
            {/*                    />*/}
            {/*                </div>*/}
            {/*                <Button type="submit" disabled={loading}>*/}
            {/*                    {loading ? 'Loading...' : 'Submit'}*/}
            {/*                </Button>*/}
            {/*                {error && <p className="text-red-500 text-sm">{error}</p>}*/}
            {/*                {responseData && (*/}
            {/*                    <div>*/}
            {/*                        <label htmlFor="response" className="block text-sm font-medium text-gray-700">*/}
            {/*                            Response Data*/}
            {/*                        </label>*/}
            {/*                        <Textarea*/}
            {/*                            id="response"*/}
            {/*                            value={responseData}*/}
            {/*                            readOnly*/}
            {/*                            className="mt-1 h-32"*/}
            {/*                        />*/}
            {/*                    </div>*/}
            {/*                )}*/}
            {/*            </form>*/}
            {/*        </PopoverContent>*/}
            {/*    </Popover>*/}
            {/*</div>*/}
        </header>
    );
}
