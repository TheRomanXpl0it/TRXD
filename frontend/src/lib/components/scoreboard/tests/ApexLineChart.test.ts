import { render } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ApexLineChart from '../ApexLineChart.svelte';

vi.mock('apexcharts', () => {
	const mock = {
		render: vi.fn(),
		updateOptions: vi.fn(),
		updateSeries: vi.fn(),
		destroy: vi.fn()
	};
	(globalThis as any).__ApexMocks = mock;
	function Ctor(this: any, _el: any, _opts: any) {
		this.render = mock.render;
		this.updateOptions = mock.updateOptions;
		this.updateSeries = mock.updateSeries;
		this.destroy = mock.destroy;
	}
	return { default: Ctor as any };
});

describe('ApexLineChart Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('instantiates chart on mount', () => {
		render(ApexLineChart, {
			props: {
				series: [{ data: [1, 2, 3] }],
				options: { chart: { id: 'x' } }
			}
		});

		expect((globalThis as any).__ApexMocks.render).toHaveBeenCalled();
	});

	it('updates options when options prop changes', async () => {
		const { rerender } = render(ApexLineChart, {
			props: {
				series: [{ data: [1, 2] }],
				options: { chart: { id: 'x' } }
			}
		});

		await rerender({
			series: [{ data: [1, 2, 3, 4] }],
			options: { chart: { id: 'y' } }
		});

		expect((globalThis as any).__ApexMocks.updateOptions).toHaveBeenCalled();
		expect((globalThis as any).__ApexMocks.updateSeries).toHaveBeenCalled();
	});

	it('destroys chart instance when component unmounts', () => {
		const { unmount } = render(ApexLineChart, {
			props: {
				series: [{ data: [1, 2, 3] }],
				options: { chart: { id: 'test' } }
			}
		});

		unmount();

		expect((globalThis as any).__ApexMocks.destroy).toHaveBeenCalled();
	});

	it('handles empty series data', () => {
		render(ApexLineChart, {
			props: {
				series: [],
				options: { chart: { id: 'empty' } }
			}
		});

		expect((globalThis as any).__ApexMocks.render).toHaveBeenCalled();
	});

	it('updates series when series prop changes', async () => {
		const { rerender } = render(ApexLineChart, {
			props: {
				series: [{ data: [1, 2, 3] }],
				options: { chart: { id: 'x' } }
			}
		});

		await rerender({
			series: [{ data: [4, 5, 6] }],
			options: { chart: { id: 'x' } }
		});

		expect((globalThis as any).__ApexMocks.updateSeries).toHaveBeenCalled();
	});
});
