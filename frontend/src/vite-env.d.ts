/// <reference types="vite/client" />

declare module 'd3-scale' {
	export function scaleTime(): any;
	export function scaleLinear(): any;
	export function scaleBand(): any;
	export function scaleOrdinal(): any;
	export function scaleLog(): any;
	export function scalePow(): any;
	export function scaleSqrt(): any;
	export function scalePoint(): any;
	export function scaleSequential(): any;
	export function scaleDiverging(): any;
	export function scaleQuantize(): any;
	export function scaleQuantile(): any;
	export function scaleThreshold(): any;
	export function scaleIdentity(): any;
	export function scaleRadial(): any;
}

declare module '@sveltejs/svelte-virtual-list' {
	import { SvelteComponent } from 'svelte';
	export default class VirtualList extends SvelteComponent<{
		items: any[];
		height?: string;
		itemHeight?: number;
		[key: string]: any;
	}> { }
}

declare global {
	const __BACKEND_URL__: string;
	const __GIT_HASH__: string;
}

interface ImportMetaEnv {
	readonly VITE_BACKEND_URL: string;
	readonly VITE_PUBLIC_URL: string;
	readonly VITE_API_BASE: string;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
