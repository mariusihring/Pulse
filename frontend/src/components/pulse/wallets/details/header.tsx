import type { Wallet } from "@/graphql/graphql";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

export default function WalletHeader({ wallet }: { wallet: Wallet }) {
	return (
		<div className="flex flex-col gap-8 pb-8">
			<div className="flex items-start justify-between">
				<div className="flex items-center gap-4">
					<div className="rounded-full bg-primary/10 p-3">
                    {/* TODO: map wallet.chain to logo from https://svgl.app/ */}
                        <Bitcoin className="h-8 w-8 text-primary" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold">{wallet.name}</h1>
                        <p className="text-sm text-muted-foreground">
                            Main Wallet • {wallet.subwallets.length} Sub-wallets
                        </p>
                    </div>
				</div>
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                        <MoreHorizontal className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-56">
                        <DropdownMenuLabel>Wallet Actions</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem>
                        <PlusCircle className="mr-2 h-4 w-4" /> Create Subwallet
                        </DropdownMenuItem>
                        <DropdownMenuItem disabled>
                        <Copy className="mr-2 h-4 w-4" /> Copy Address
                        </DropdownMenuItem>
                        <DropdownMenuItem disabled>
                        <QrCode className="mr-2 h-4 w-4" /> Show QR Code
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem className="text-red-600">
                        Remove Wallet
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                    </DropdownMenu>
			</div>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div className="space-y-2">
          <p className="text-sm text-muted-foreground">Total Balance</p>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold">
              $1.000.000
            </span>
            <span className="text-lg text-muted-foreground">
              69 BTC
            </span>
          </div>
        </div>
        <div className="space-y-2">
          <p className="text-sm text-muted-foreground">24h Change</p>
          <div className="flex items-center gap-2">
            {12 >= 0 ? (
              <ArrowUpRight className="h-4 w-4 text-green-500" />
            ) : (
              <ArrowDownRight className="h-4 w-4 text-red-500" />
            )}
            <span className={`text-2xl font-bold ${
              12 >= 0 ? "text-green-500" : "text-red-500"
            }`}>
              {12 >= 0 ? "+" : ""}{12}%
            </span>
          </div>
        </div>
      </div>
      <div className="flex flex-wrap gap-2">
        <Button disabled>
          <Send className="mr-2 h-4 w-4" /> Send
        </Button>
        <Button disabled>
          <Download className="mr-2 h-4 w-4" /> Receive
        </Button>
        <Button variant="outline" disabled>
          <Repeat className="mr-2 h-4 w-4" /> Swap
        </Button>
      </div>
		</div>
	);
}

import type { SVGProps } from "react";
import { Button } from "@/components/ui/button";
import { MoreHorizontal, Copy, QrCode, ArrowDownRight, ArrowUpRight, Download, Repeat, Send, PlusCircle } from "lucide-react";
const Bitcoin = (props: SVGProps<SVGSVGElement>) => (
	// biome-ignore lint/a11y/noSvgWithoutTitle: <explanation>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		xmlnsXlink="http://www.w3.org/1999/xlink"
		width="1em"
		height="1em"
		viewBox="0 0 32 32"
		{...props}
	>
		<defs>
			<linearGradient id="btc-c" x1="50%" x2="50%" y1="0%" y2="100%">
				<stop offset="0%" stopColor="#FFF" stopOpacity={0.5} />
				<stop offset="100%" stopOpacity={0.5} />
			</linearGradient>
			<circle id="btc-b" cx={16} cy={15} r={15} />
			<filter
				id="btc-a"
				width="111.7%"
				height="111.7%"
				x="-5.8%"
				y="-4.2%"
				filterUnits="objectBoundingBox"
			>
				<feOffset dy={0.5} in="SourceAlpha" result="shadowOffsetOuter1" />
				<feGaussianBlur
					in="shadowOffsetOuter1"
					result="shadowBlurOuter1"
					stdDeviation={0.5}
				/>
				<feComposite
					in="shadowBlurOuter1"
					in2="SourceAlpha"
					operator="out"
					result="shadowBlurOuter1"
				/>
				<feColorMatrix
					in="shadowBlurOuter1"
					values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.199473505 0"
				/>
			</filter>
			<path
				id="btc-e"
				d="M23.1889526,13.0201846 C23.5025526,10.9239385 21.9064911,9.79704615 19.7240911,9.04529231 L20.4320295,6.20566154 L18.7035372,5.77489231 L18.0143065,8.53969231 C17.5599065,8.42646154 17.0931988,8.31963077 16.6294449,8.21378462 L17.3235988,5.43076923 L15.5960911,5 L14.8876603,7.83864615 C14.5115372,7.75298462 14.1423065,7.66830769 13.7839065,7.5792 L13.7858757,7.57033846 L11.4021218,6.97513846 L10.9423065,8.82129231 C10.9423065,8.82129231 12.224768,9.1152 12.1976911,9.13341538 C12.8977526,9.30818462 13.0242757,9.77144615 13.0031065,10.1387077 L12.1967065,13.3736615 C12.2449526,13.3859692 12.3074757,13.4036923 12.3763988,13.4312615 C12.3187988,13.4169846 12.2572603,13.4012308 12.1937526,13.3859692 L11.0634142,17.9176615 C10.9777526,18.1303385 10.7606449,18.4493538 10.2712911,18.3282462 C10.2885218,18.3533538 9.01492185,18.0146462 9.01492185,18.0146462 L8.15682954,19.9932308 L10.4061834,20.5539692 C10.8246449,20.6588308 11.2347372,20.7686154 11.6384295,20.872 L10.9231065,23.7441231 L12.6496295,24.1748923 L13.3580603,21.3332923 C13.8296911,21.4612923 14.2875372,21.5794462 14.7355372,21.6907077 L14.029568,24.5190154 L15.7580603,24.9497846 L16.4733834,22.0830769 C19.4208295,22.6408615 21.6371988,22.4158769 22.5701218,19.7500308 C23.3218757,17.6035692 22.5327065,16.3654154 20.9819372,15.5580308 C22.1112911,15.2976 22.9619988,14.5547077 23.1889526,13.0201846 L23.1889526,13.0201846 Z M19.2396603,18.5581538 C18.7055065,20.7046154 15.0914757,19.5442462 13.9197834,19.2532923 L14.8689526,15.4482462 C16.0406449,15.7406769 19.7979372,16.3196308 19.2396603,18.5581538 Z M19.7743065,12.9891692 C19.2869218,14.9416615 16.2789218,13.9496615 15.303168,13.7064615 L16.1637218,10.2553846 C17.1394757,10.4985846 20.2818757,10.9524923 19.7743065,12.9891692 Z"
			/>
			<filter
				id="btc-d"
				width="123.2%"
				height="117.5%"
				x="-11.6%"
				y="-6.3%"
				filterUnits="objectBoundingBox"
			>
				<feOffset dy={0.5} in="SourceAlpha" result="shadowOffsetOuter1" />
				<feGaussianBlur
					in="shadowOffsetOuter1"
					result="shadowBlurOuter1"
					stdDeviation={0.5}
				/>
				<feColorMatrix
					in="shadowBlurOuter1"
					values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.204257246 0"
				/>
			</filter>
		</defs>
		<g fill="none" fillRule="evenodd">
			<use fill="#000" filter="url(#btc-a)" xlinkHref="#btc-b" />
			<use fill="#F7931A" xlinkHref="#btc-b" />
			<use
				fill="url(#btc-c)"
				style={{
					mixBlendMode: "soft-light",
				}}
				xlinkHref="#btc-b"
			/>
			<circle cx={16} cy={15} r={14.5} stroke="#000" strokeOpacity={0.097} />
			<g fillRule="nonzero">
				<use fill="#000" filter="url(#btc-d)" xlinkHref="#btc-e" />
				<use fill="#FFF" fillRule="evenodd" xlinkHref="#btc-e" />
			</g>
		</g>
	</svg>
);
