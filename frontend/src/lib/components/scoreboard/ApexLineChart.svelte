<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import ApexCharts from "apexcharts";

  const p: any = $props();

  let el: HTMLDivElement | null = null;
  let chart: any = null;
  let ro: ResizeObserver | null = null;

  function safeSeries(input: any): any[] {
    if (Array.isArray(input)) return input;
    if (input && typeof input === "object") return [input];
    return [];
  }

  function buildOptions(): any {
    const options = p.options ?? {};
    const o: any = {
      ...options,
      chart: {
        ...(options.chart ?? {}),
        type: p.type ?? "line",
        height: p.height ?? 360,
        width: p.width ?? "100%"
      }
    };
    // ensure we don't pass a non-array series inside options
    if (o.series && !Array.isArray(o.series)) delete o.series;
    return o;
  }

  function nudgeResize() {
    if (!chart) return;
    try { chart.updateOptions({}, true, true); } catch {}
  }

  onMount(() => {
    if (!el) return;
    const initial = { ...buildOptions(), series: safeSeries(p.series) };
    chart = new (ApexCharts as any)(el, initial);
    chart.render();

    if (typeof ResizeObserver !== "undefined") {
      ro = new ResizeObserver(() => nudgeResize());
      ro.observe(el);
    } else {
      setTimeout(() => window.dispatchEvent(new Event("resize")), 0);
    }
  });

  // Update options when option-like props change
  $effect(() => {
    if (!chart) return;
    void p.options; void p.type; void p.height; void p.width;
    try { chart.updateOptions(buildOptions(), false, true); } catch (e) {
      console.error("[ApexLineChart] updateOptions:", e);
    }
  });

  // Update series separately (keeps viewport)
  $effect(() => {
    if (!chart) return;
    void p.series;
    try { chart.updateSeries(safeSeries(p.series), true); } catch (e) {
      console.error("[ApexLineChart] updateSeries:", e, p.series);
    }
  });

  onDestroy(() => {
    try { ro?.disconnect(); } catch {}
    ro = null;
    try { chart?.destroy(); } finally { chart = null; }
  });
</script>

<div bind:this={el} style="width: {typeof (p.width ?? '100%') === 'number' ? `${p.width}px` : (p.width ?? '100%')};"></div>
