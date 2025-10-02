import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { Line } from "react-chartjs-2";
import { useEffect, useState } from "react";
import { getScoreboard } from "@/lib/backend-interaction";
import type { ScoreboardEntry } from "@/components/Scores";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
);

export const options = {
  responsive: true,
  interaction: {
    mode: "index" as const,
    intersect: false,
  },
  stacked: false,
  plugins: {
    title: {
      display: true,
      text: "Scoreboard Chart",
    },
  },
  scales: {
    y: {
      type: "linear" as const,
      display: true,
      position: "left" as const,
    },
    y1: {
      type: "linear" as const,
      display: true,
      position: "right" as const,
      grid: {
        drawOnChartArea: false,
      },
    },
  },
};

function convertScoreboardEntryToGraphData(data: ScoreboardEntry[]) {
  const graphData = {
    labels: data.map((entry) => entry.name),
    datasets: [
      {
        label: "Scores",
        data: data.map((entry) => entry.score),
        borderColor: "rgb(75, 192, 192)",
        backgroundColor: "rgba(75, 192, 192, 0.2)",
        yAxisID: "y",
      },
    ],
  };
  return graphData;
}

export function Chart() {
  const [scores, setScores] = useState<ScoreboardEntry[]>([]);

  useEffect(() => {
    let mounted = true;
    (async () => {
      const result = await getScoreboard();
      if (mounted && Array.isArray(result)) {
        setScores(result as ScoreboardEntry[]);
      } else if (mounted) {
        setScores([]);
      }
    })();
    return () => {
      mounted = false;
    };
  }, []);

  const graphData = convertScoreboardEntryToGraphData(scores);
  return (
    <>
      <Line options={options} data={graphData} />
    </>
  );
}
