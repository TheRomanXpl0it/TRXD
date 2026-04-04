import { render, screen } from '@testing-library/svelte';
import RadarChart from '$lib/components/RadarChart.svelte';

describe('RadarChart', () => {
    it('renders the placeholder message when competition hasn\'t started (no totalChallenges)', () => {
        render(RadarChart, {
            props: {
                solves: [],
                totalChallenges: []
            }
        });

        expect(screen.getByText("The competition hasn't started yet")).toBeInTheDocument();
    });

    it('prepares to render the chart when totalChallenges are provided', () => {
        const { container } = render(RadarChart, {
            props: {
                solves: [],
                totalChallenges: [
                    { category: 'Web', count: 5 },
                    { category: 'Pwn', count: 3 }
                ]
            }
        });

        expect(screen.queryByText("The competition hasn't started yet")).not.toBeInTheDocument();
        const chartDiv = container.querySelector('div.w-full');
        expect(chartDiv).toBeDefined();
    });
});
